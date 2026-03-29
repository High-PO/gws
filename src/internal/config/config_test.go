package config

import (
	"testing"

	"gws/internal/credential"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// --- Unit Tests ---

func TestGetMFASerial_ErrorWhenNoSerial(t *testing.T) {
	store := credential.NewMockStore()
	mgr := &Manager{Store: store}

	_, err := mgr.GetMFASerial("default")
	if err == nil {
		t.Fatal("expected error when no serial is stored, got nil")
	}
}

func TestSaveThenGetMFASerial(t *testing.T) {
	store := credential.NewMockStore()
	mgr := &Manager{Store: store}

	if err := mgr.SaveMFASerial("default", "arn:aws:iam::123456789012:mfa/user"); err != nil {
		t.Fatalf("SaveMFASerial failed: %v", err)
	}

	got, err := mgr.GetMFASerial("default")
	if err != nil {
		t.Fatalf("GetMFASerial failed: %v", err)
	}
	if got != "arn:aws:iam::123456789012:mfa/user" {
		t.Errorf("GetMFASerial = %q, want %q", got, "arn:aws:iam::123456789012:mfa/user")
	}
}

func TestOverwriteMFASerial(t *testing.T) {
	store := credential.NewMockStore()
	mgr := &Manager{Store: store}

	mgr.SaveMFASerial("prod", "old-serial")
	mgr.SaveMFASerial("prod", "new-serial")

	got, err := mgr.GetMFASerial("prod")
	if err != nil {
		t.Fatalf("GetMFASerial failed: %v", err)
	}
	if got != "new-serial" {
		t.Errorf("GetMFASerial = %q, want %q", got, "new-serial")
	}
}

// --- Property-Based Tests ---

// Feature: gws-refactoring, Property 2: MFA 시리얼 저장/조회 라운드트립
// For any profile name and MFA serial string, saving then getting should return the same value.
// **Validates: Requirements 2.1, 2.2**
func TestProperty2_MFASerialRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	properties.Property("SaveMFASerial then GetMFASerial returns the same value", prop.ForAll(
		func(profile string, serial string) bool {
			store := credential.NewMockStore()
			mgr := &Manager{Store: store}

			if err := mgr.SaveMFASerial(profile, serial); err != nil {
				return false
			}

			got, err := mgr.GetMFASerial(profile)
			if err != nil {
				return false
			}

			return got == serial
		},
		gen.AlphaString().WithLabel("profile"),
		gen.AlphaString().WithLabel("serial"),
	))

	properties.TestingRun(t)
}

// Feature: gws-refactoring, Property 3: 프로필별 MFA 시리얼 독립성
// For any two different profile names with different serials, saving to one profile should not affect the other.
// **Validates: Requirements 2.5**
func TestProperty3_ProfileIndependence(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	properties.Property("saving to one profile does not affect another profile", prop.ForAll(
		func(profile1, profile2, serial1, serial2 string) bool {
			if profile1 == profile2 {
				return true // skip when profiles are the same
			}

			store := credential.NewMockStore()
			mgr := &Manager{Store: store}

			// Save serial for both profiles
			if err := mgr.SaveMFASerial(profile1, serial1); err != nil {
				return false
			}
			if err := mgr.SaveMFASerial(profile2, serial2); err != nil {
				return false
			}

			// Verify profile1 still has its own serial
			got1, err := mgr.GetMFASerial(profile1)
			if err != nil {
				return false
			}
			if got1 != serial1 {
				return false
			}

			// Verify profile2 has its own serial
			got2, err := mgr.GetMFASerial(profile2)
			if err != nil {
				return false
			}
			if got2 != serial2 {
				return false
			}

			// Now overwrite profile2 and verify profile1 is unaffected
			if err := mgr.SaveMFASerial(profile2, "changed"); err != nil {
				return false
			}

			got1After, err := mgr.GetMFASerial(profile1)
			if err != nil {
				return false
			}

			return got1After == serial1
		},
		gen.AlphaString().WithLabel("profile1"),
		gen.AlphaString().WithLabel("profile2"),
		gen.AlphaString().WithLabel("serial1"),
		gen.AlphaString().WithLabel("serial2"),
	))

	properties.TestingRun(t)
}
