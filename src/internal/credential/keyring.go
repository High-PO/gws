package credential

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

// KeyringStore implements the Store interface using OS-native secure storage
// via go-keyring (macOS Keychain, Linux Secret Service, Windows Credential Manager).
type KeyringStore struct{}

// Get retrieves a value from the OS-native secure storage.
func (k *KeyringStore) Get(service, key string) (string, error) {
	val, err := keyring.Get(service, key)
	if err != nil {
		return "", fmt.Errorf("보안 저장소 조회 실패: %w", err)
	}
	return val, nil
}

// Set stores a value in the OS-native secure storage.
func (k *KeyringStore) Set(service, key, value string) error {
	if err := keyring.Set(service, key, value); err != nil {
		return fmt.Errorf("보안 저장소 저장 실패: %w", err)
	}
	return nil
}

// Delete removes a value from the OS-native secure storage.
func (k *KeyringStore) Delete(service, key string) error {
	if err := keyring.Delete(service, key); err != nil {
		return fmt.Errorf("보안 저장소 삭제 실패: %w", err)
	}
	return nil
}

// NewStore returns a Store implementation backed by the OS-native secure storage.
func NewStore() Store {
	return &KeyringStore{}
}
