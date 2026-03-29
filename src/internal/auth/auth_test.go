package auth

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// --- Mock Implementations ---

// mockSTSClient implements STSClient with configurable return values.
type mockSTSClient struct {
	creds *SessionCredential
	err   error
}

func (m *mockSTSClient) GetSessionToken(profile, serialNumber, tokenCode string) (*SessionCredential, error) {
	return m.creds, m.err
}

// mockConfigProvider implements ConfigProvider with in-memory storage.
type mockConfigProvider struct {
	serials map[string]string
}

func newMockConfigProvider() *mockConfigProvider {
	return &mockConfigProvider{serials: make(map[string]string)}
}

func (m *mockConfigProvider) GetMFASerial(profile string) (string, error) {
	s, ok := m.serials[profile]
	if !ok {
		return "", fmt.Errorf("no serial for profile %q", profile)
	}
	return s, nil
}

func (m *mockConfigProvider) SaveMFASerial(profile, serial string) error {
	m.serials[profile] = serial
	return nil
}

// mockShellLauncher implements ShellLauncher that records calls.
type mockShellLauncher struct {
	launched bool
	creds    *SessionCredential
	profile  string
	err      error
}

func (m *mockShellLauncher) Launch(creds *SessionCredential, profile string) error {
	m.launched = true
	m.creds = creds
	m.profile = profile
	return m.err
}

// --- Unit Tests ---

func TestAuthenticate_Success(t *testing.T) {
	creds := &SessionCredential{
		AccessKeyID:     "ABCDEFGHIJKLMNOPQRST",
		SecretAccessKey: "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcd",
		SessionToken:    "test-session-token-value",
		Expiration:      "2024-01-01T12:00:00Z",
	}

	sts := &mockSTSClient{creds: creds}
	cfg := newMockConfigProvider()
	cfg.serials["default"] = "arn:aws:iam::000000000000:mfa/testdevice"
	shell := &mockShellLauncher{}

	auth := &Authenticator{STS: sts, Config: cfg, Shell: shell}
	err := auth.Authenticate("default", "123456")
	if err != nil {
		t.Fatalf("Authenticate failed: %v", err)
	}
	if !shell.launched {
		t.Error("expected shell to be launched")
	}
	if shell.creds.AccessKeyID != creds.AccessKeyID {
		t.Errorf("shell creds AccessKeyID = %q, want %q", shell.creds.AccessKeyID, creds.AccessKeyID)
	}
}

func TestAuthenticate_STSFailure(t *testing.T) {
	sts := &mockSTSClient{err: fmt.Errorf("ExpiredTokenException: token expired")}
	cfg := newMockConfigProvider()
	cfg.serials["default"] = "arn:aws:iam::000000000000:mfa/testdevice"
	shell := &mockShellLauncher{}

	auth := &Authenticator{STS: sts, Config: cfg, Shell: shell}
	err := auth.Authenticate("default", "123456")
	if err == nil {
		t.Fatal("expected error from STS failure, got nil")
	}
	if !strings.Contains(err.Error(), "ExpiredTokenException") {
		t.Errorf("error %q should contain STS cause", err.Error())
	}
	if shell.launched {
		t.Error("shell should not be launched on STS failure")
	}
}

func TestAuthenticate_ShellLaunchFailure(t *testing.T) {
	creds := &SessionCredential{
		AccessKeyID:     "TESTKEY",
		SecretAccessKey: "TESTSECRET",
		SessionToken:    "TESTTOKEN",
		Expiration:      "2024-01-01T12:00:00Z",
	}

	sts := &mockSTSClient{creds: creds}
	cfg := newMockConfigProvider()
	cfg.serials["default"] = "arn:aws:iam::000000000000:mfa/testdevice"
	shell := &mockShellLauncher{err: fmt.Errorf("shell not found")}

	auth := &Authenticator{STS: sts, Config: cfg, Shell: shell}
	err := auth.Authenticate("default", "123456")
	if err == nil {
		t.Fatal("expected error from shell launch failure, got nil")
	}
	if !strings.Contains(err.Error(), "shell not found") {
		t.Errorf("error %q should contain shell cause", err.Error())
	}
}

// --- Property-Based Tests ---

// Feature: gws-refactoring, Property 6: error message contains cause
// For any error cause string, when STS fails with that cause, the returned error
// message should contain the original cause string. Also test JSON parse failure case.
// **Validates: Requirements 4.1, 4.2**
func TestProperty6_ErrorMessageContainsCause(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	properties.Property("STS failure error contains cause string", prop.ForAll(
		func(cause string) bool {
			sts := &mockSTSClient{err: fmt.Errorf("some cause: %s", cause)}
			cfg := newMockConfigProvider()
			cfg.serials["default"] = "arn:aws:iam::000000000000:mfa/testdevice"
			shell := &mockShellLauncher{}

			auth := &Authenticator{STS: sts, Config: cfg, Shell: shell}
			err := auth.Authenticate("default", "123456")
			if err == nil {
				return false
			}
			return strings.Contains(err.Error(), cause)
		},
		gen.AlphaString().WithLabel("cause"),
	))

	properties.Property("JSON parse failure error contains cause", prop.ForAll(
		func(badJSON string) bool {
			_, err := ParseSTSResponse([]byte(badJSON))
			if err == nil {
				// If it happens to be valid JSON with Credentials, that's fine
				return true
			}
			// The error should contain some indication of the parse failure
			return strings.Contains(err.Error(), "파싱 실패") || strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "cannot") || strings.Contains(err.Error(), "unexpected")
		},
		gen.AnyString().SuchThat(func(s string) bool {
			// Filter to strings that are NOT valid STS JSON
			var resp stsResponse
			return json.Unmarshal([]byte(s), &resp) != nil
		}).WithLabel("badJSON"),
	))

	properties.TestingRun(t)
}

// Feature: gws-refactoring, Property 8: STS response JSON round-trip
// For any valid SessionCredential (non-empty fields), serializing to the STS JSON
// format {"Credentials": {...}} and parsing back with ParseSTSResponse should return
// the same credential values.
// **Validates: Requirements 4.2**
func TestProperty8_STSResponseJSONRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	nonEmptyAlpha := gen.AlphaString().SuchThat(func(s string) bool {
		return len(s) > 0
	})

	properties.Property("JSON serialize then ParseSTSResponse returns same credentials", prop.ForAll(
		func(accessKeyID, secretAccessKey, sessionToken, expiration string) bool {
			// Serialize to STS JSON format
			wrapper := stsResponse{
				Credentials: SessionCredential{
					AccessKeyID:     accessKeyID,
					SecretAccessKey: secretAccessKey,
					SessionToken:    sessionToken,
					Expiration:      expiration,
				},
			}
			data, err := json.Marshal(wrapper)
			if err != nil {
				return false
			}

			// Parse back
			got, err := ParseSTSResponse(data)
			if err != nil {
				return false
			}

			return got.AccessKeyID == accessKeyID &&
				got.SecretAccessKey == secretAccessKey &&
				got.SessionToken == sessionToken &&
				got.Expiration == expiration
		},
		nonEmptyAlpha.WithLabel("AccessKeyID"),
		nonEmptyAlpha.WithLabel("SecretAccessKey"),
		nonEmptyAlpha.WithLabel("SessionToken"),
		nonEmptyAlpha.WithLabel("Expiration"),
	))

	properties.TestingRun(t)
}
