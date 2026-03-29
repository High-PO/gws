package auth

import "fmt"

// ConfigProvider abstracts config operations (MFA serial get/save).
type ConfigProvider interface {
	GetMFASerial(profile string) (string, error)
	SaveMFASerial(profile, serial string) error
}

// ShellLauncher abstracts shell execution.
type ShellLauncher interface {
	Launch(creds *SessionCredential, profile string) error
}

// Authenticator performs MFA authentication.
type Authenticator struct {
	STS    STSClient
	Config ConfigProvider
	Shell  ShellLauncher
}

// Authenticate performs the full MFA authentication flow:
// 1. Get MFA serial from config (or prompt user if not found)
// 2. Call STS get-session-token
// 3. Print expiration time
// 4. Launch shell with credentials
func (a *Authenticator) Authenticate(profile, token string) error {
	serialNumber, err := a.Config.GetMFASerial(profile)
	if err != nil {
		fmt.Print("Enter MFA serial number (arn:aws:iam::<account-id>:mfa/<device>): ")
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			return fmt.Errorf("MFA 시리얼 번호 읽기 실패: %w", err)
		}
		if input == "" {
			return fmt.Errorf("MFA serial number cannot be empty")
		}
		serialNumber = input
		if err := a.Config.SaveMFASerial(profile, serialNumber); err != nil {
			return fmt.Errorf("MFA 시리얼 저장 실패: %w", err)
		}
	}

	creds, err := a.STS.GetSessionToken(profile, serialNumber, token)
	if err != nil {
		return fmt.Errorf("STS API 호출 실패: %w", err)
	}

	fmt.Printf("AWS session credentials set. Expires at: %s\n", creds.Expiration)

	return a.Shell.Launch(creds, profile)
}
