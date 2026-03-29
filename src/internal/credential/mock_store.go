package credential

import "fmt"

// MockStore implements the Store interface using in-memory map storage.
// It is exported so other packages can reuse it for testing.
type MockStore struct {
	data map[string]map[string]string
}

// NewMockStore creates a new MockStore with an initialized internal map.
func NewMockStore() *MockStore {
	return &MockStore{
		data: make(map[string]map[string]string),
	}
}

func (m *MockStore) Get(service, key string) (string, error) {
	svc, ok := m.data[service]
	if !ok {
		return "", fmt.Errorf("secret not found for service %q, key %q", service, key)
	}
	val, ok := svc[key]
	if !ok {
		return "", fmt.Errorf("secret not found for service %q, key %q", service, key)
	}
	return val, nil
}

func (m *MockStore) Set(service, key, value string) error {
	if m.data[service] == nil {
		m.data[service] = make(map[string]string)
	}
	m.data[service][key] = value
	return nil
}

func (m *MockStore) Delete(service, key string) error {
	svc, ok := m.data[service]
	if !ok {
		return fmt.Errorf("secret not found for service %q, key %q", service, key)
	}
	if _, ok := svc[key]; !ok {
		return fmt.Errorf("secret not found for service %q, key %q", service, key)
	}
	delete(svc, key)
	return nil
}

// Verify MockStore implements Store at compile time.
var _ Store = (*MockStore)(nil)
