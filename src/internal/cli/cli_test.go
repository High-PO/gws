package cli

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// --- Unit Tests ---

func TestParse_ZeroArgs_CmdUsage(t *testing.T) {
	cmd, err := Parse([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Type != CmdUsage {
		t.Errorf("Parse([]) Type = %v, want CmdUsage", cmd.Type)
	}
}

func TestParse_Help(t *testing.T) {
	cmd, err := Parse([]string{"help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Type != CmdHelp {
		t.Errorf("Parse([help]) Type = %v, want CmdHelp", cmd.Type)
	}
}

func TestParse_Version(t *testing.T) {
	cmd, err := Parse([]string{"--version"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Type != CmdVersion {
		t.Errorf("Parse([--version]) Type = %v, want CmdVersion", cmd.Type)
	}
}

func TestParse_ValidToken_DefaultProfile(t *testing.T) {
	cmd, err := Parse([]string{"123456"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Type != CmdAuth {
		t.Errorf("Type = %v, want CmdAuth", cmd.Type)
	}
	if cmd.Profile != "default" {
		t.Errorf("Profile = %q, want %q", cmd.Profile, "default")
	}
	if cmd.Token != "123456" {
		t.Errorf("Token = %q, want %q", cmd.Token, "123456")
	}
}

func TestParse_ProfileAndToken(t *testing.T) {
	cmd, err := Parse([]string{"production", "654321"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Type != CmdAuth {
		t.Errorf("Type = %v, want CmdAuth", cmd.Type)
	}
	if cmd.Profile != "production" {
		t.Errorf("Profile = %q, want %q", cmd.Profile, "production")
	}
	if cmd.Token != "654321" {
		t.Errorf("Token = %q, want %q", cmd.Token, "654321")
	}
}

func TestParse_InvalidSingleArg(t *testing.T) {
	_, err := Parse([]string{"notacommand"})
	if err == nil {
		t.Fatal("expected error for invalid single arg, got nil")
	}
}

func TestParse_ProfileAndInvalidToken(t *testing.T) {
	_, err := Parse([]string{"myprofile", "abc"})
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

func TestParse_TooManyArgs(t *testing.T) {
	_, err := Parse([]string{"a", "b", "c"})
	if err == nil {
		t.Fatal("expected error for 3+ args, got nil")
	}
}

func TestValidateMFAToken_Valid(t *testing.T) {
	validTokens := []string{"000000", "123456", "999999", "001234"}
	for _, tok := range validTokens {
		if err := ValidateMFAToken(tok); err != nil {
			t.Errorf("ValidateMFAToken(%q) returned error: %v", tok, err)
		}
	}
}

func TestValidateMFAToken_Invalid(t *testing.T) {
	invalidTokens := []string{"", "12345", "1234567", "abcdef", "12345a", "12 345"}
	for _, tok := range invalidTokens {
		if err := ValidateMFAToken(tok); err == nil {
			t.Errorf("ValidateMFAToken(%q) expected error, got nil", tok)
		}
	}
}

// --- Property-Based Tests ---

var sixDigitPattern = regexp.MustCompile(`^[0-9]{6}$`)

// Feature: gws-refactoring, Property 4: CLI parsing correctness
// For any valid profile name (alphanumeric) and 6-digit numeric token:
// - [token] format parses to profile="default" and the given token
// - [profile, token] format parses to the given profile and token
// **Validates: Requirements 5.1, 5.2**
func TestProperty4_CLIParsingCorrectness(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	tokenGen := gen.IntRange(0, 999999).Map(func(n int) string {
		return fmt.Sprintf("%06d", n)
	})

	profileGen := gen.AlphaString().SuchThat(func(s string) bool {
		return len(s) > 0
	})

	properties.Property("single token arg parses to default profile", prop.ForAll(
		func(token string) bool {
			cmd, err := Parse([]string{token})
			if err != nil {
				return false
			}
			return cmd.Type == CmdAuth && cmd.Profile == "default" && cmd.Token == token
		},
		tokenGen.WithLabel("token"),
	))

	properties.Property("profile + token args parse correctly", prop.ForAll(
		func(profile, token string) bool {
			cmd, err := Parse([]string{profile, token})
			if err != nil {
				return false
			}
			return cmd.Type == CmdAuth && cmd.Profile == profile && cmd.Token == token
		},
		profileGen.WithLabel("profile"),
		tokenGen.WithLabel("token"),
	))

	properties.TestingRun(t)
}

// Feature: gws-refactoring, Property 5: MFA token validation
// For any string, ValidateMFAToken returns nil if and only if the string matches ^[0-9]{6}$
// **Validates: Requirements 5.5, 5.6**
func TestProperty5_MFATokenValidation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	properties.Property("ValidateMFAToken returns nil iff string is exactly 6 digits", prop.ForAll(
		func(s string) bool {
			err := ValidateMFAToken(s)
			isValid := sixDigitPattern.MatchString(s)
			if isValid {
				return err == nil
			}
			return err != nil
		},
		gen.AnyString().WithLabel("input"),
	))

	properties.TestingRun(t)
}

// Feature: gws-refactoring, Property 7: Invalid argument rejection
// For any args with 3+ elements, Parse returns an error.
// For 1-arg that is not "help", "--version", or a 6-digit number, Parse returns error.
// **Validates: Requirements 4.3**
func TestProperty7_InvalidArgumentRejection(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	// 3+ args always produce an error
	properties.Property("3+ args always returns error", prop.ForAll(
		func(args []string) bool {
			_, err := Parse(args)
			return err != nil
		},
		gen.SliceOfN(3, gen.AnyString()).SuchThat(func(s []string) bool {
			return len(s) >= 3
		}).WithLabel("args"),
	))

	// 1 invalid arg (not "help", "--version", or 6-digit number) returns error
	invalidSingleArgGen := gen.AlphaString().SuchThat(func(s string) bool {
		if s == "help" || s == "--version" {
			return false
		}
		if sixDigitPattern.MatchString(s) {
			return false
		}
		return true
	})

	properties.Property("single invalid arg returns error", prop.ForAll(
		func(arg string) bool {
			_, err := Parse([]string{arg})
			return err != nil
		},
		invalidSingleArgGen.WithLabel("invalidArg"),
	))

	properties.TestingRun(t)
}
