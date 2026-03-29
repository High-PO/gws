package cli

import (
	"fmt"
	"io"
)

func printHelpEn(w io.Writer) {
	fmt.Fprintln(w, "GWS - Go + AWS CLI v2")
	fmt.Fprintln(w, "=====================")
	fmt.Fprintln(w, "A tool to easily use AWS CLI v2 as an MFA-enabled IAM user.")
	fmt.Fprintln(w, "Supports macOS, Linux, and Windows.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "USAGE:")
	fmt.Fprintln(w, "  gws <mfa-token>                  # Use default AWS profile with MFA token")
	fmt.Fprintln(w, "  gws <profile> <mfa-token>        # Use specific AWS profile with MFA token")
	fmt.Fprintln(w, "  gws help                         # Show this help message")
	fmt.Fprintln(w, "  gws --version                    # Show version information")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "EXAMPLES:")
	fmt.Fprintln(w, "  1. Using default profile:")
	fmt.Fprintln(w, "     $ gws 123456")
	fmt.Fprintln(w, "     This will authenticate using your default AWS profile with MFA token 123456")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "  2. Using a specific profile:")
	fmt.Fprintln(w, "     $ gws production 654321")
	fmt.Fprintln(w, "     This will authenticate using the 'production' AWS profile with MFA token 654321")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "ARGUMENTS:")
	fmt.Fprintln(w, "  <profile>      AWS profile name configured in ~/.aws/config")
	fmt.Fprintln(w, "  <mfa-token>    6-digit code from your MFA device (e.g., Google Authenticator)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "NOTES:")
	fmt.Fprintln(w, "  - MFA serial numbers are stored securely in OS secure storage per profile")
	fmt.Fprintln(w, "  - On first use with a profile, you'll be prompted for the MFA serial number")
	fmt.Fprintln(w, "  - The tool creates a new shell session with temporary AWS credentials")
	fmt.Fprintln(w, "  - Credentials are valid for 12 hours by default (AWS STS limitation)")
}

func printUsageEn(w io.Writer) {
	fmt.Fprintln(w, "Usage: gws [profile] <mfa-token>")
	fmt.Fprintln(w, "       gws help")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Run 'gws help' for detailed usage information.")
}
