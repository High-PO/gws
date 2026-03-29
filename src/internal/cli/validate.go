package cli

import "fmt"

// ValidateMFAToken validates that the given token is exactly 6 digits.
func ValidateMFAToken(token string) error {
	if !isSixDigitNumber(token) {
		return fmt.Errorf("invalid MFA token format: %q (must be exactly 6 digits)", token)
	}
	return nil
}
