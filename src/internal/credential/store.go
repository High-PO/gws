package credential

// Store is the interface for secure credential storage.
// It abstracts OS-native secure storage (macOS Keychain, Linux Secret Service, Windows Credential Manager).
type Store interface {
	Get(service, key string) (string, error)
	Set(service, key, value string) error
	Delete(service, key string) error
}
