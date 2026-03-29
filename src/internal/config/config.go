package config

import (
	"fmt"

	"gws/internal/credential"
)

const serviceName = "gws"

// Manager manages GWS configuration using a credential Store.
type Manager struct {
	Store credential.Store
}

// GetMFASerial retrieves the MFA serial for the given profile.
// Key format: "mfa-serial-{profile}"
func (m *Manager) GetMFASerial(profile string) (string, error) {
	val, err := m.Store.Get(serviceName, "mfa-serial-"+profile)
	if err != nil {
		return "", fmt.Errorf("MFA 시리얼 조회 실패: %w", err)
	}
	return val, nil
}

// SaveMFASerial saves the MFA serial for the given profile.
func (m *Manager) SaveMFASerial(profile, serial string) error {
	if err := m.Store.Set(serviceName, "mfa-serial-"+profile, serial); err != nil {
		return fmt.Errorf("MFA 시리얼 저장 실패: %w", err)
	}
	return nil
}
