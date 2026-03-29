package shell

import (
	"fmt"
	"testing"

	"gws/internal/auth"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Feature: shell-session-switch, Property 1: 셸 환경 변수 설정
// For any SessionCredential and profile name, BuildCommand should set AWS_ACCESS_KEY_ID,
// AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN, and GWS_SESSION in the command's environment.
// **Validates: Requirements 1.1, 2.2**
func TestProperty1_ShellEnvironmentVariables(t *testing.T) {
	params := gopter.DefaultTestParametersWithSeed(1)
	params.MinSuccessfulTests = 100
	properties := gopter.NewProperties(params)

	nonEmptyAlpha := gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 })

	properties.Property("BuildCommand sets AWS credential and GWS_SESSION env vars matching input", prop.ForAll(
		func(accessKey, secretKey, sessionToken, profile string) bool {
			creds := &auth.SessionCredential{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
				SessionToken:    sessionToken,
			}

			launcher := &Launcher{}
			cmd := launcher.BuildCommand(creds, profile)

			expectedAccessKey := fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", accessKey)
			expectedSecretKey := fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", secretKey)
			expectedToken := fmt.Sprintf("AWS_SESSION_TOKEN=%s", sessionToken)
			expectedSession := fmt.Sprintf("GWS_SESSION=%s", profile)

			foundAccessKey := false
			foundSecretKey := false
			foundToken := false
			foundSession := false

			for _, env := range cmd.Env {
				switch env {
				case expectedAccessKey:
					foundAccessKey = true
				case expectedSecretKey:
					foundSecretKey = true
				case expectedToken:
					foundToken = true
				case expectedSession:
					foundSession = true
				}
			}

			return foundAccessKey && foundSecretKey && foundToken && foundSession
		},
		nonEmptyAlpha,
		nonEmptyAlpha,
		nonEmptyAlpha,
		nonEmptyAlpha,
	))

	properties.TestingRun(t)
}
