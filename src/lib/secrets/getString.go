package secrets

import "fmt"

// Get retrieves a secret string for the provided key
func (m *manager) GetString(key string) (string, error) {
	sec, err := m.Get(key)
	if err != nil {
		return "", err
	}
	secStr, ok := sec.(string)
	if secStr == "" || !ok {
		return "", fmt.Errorf("Secret is not available")
	}
	return secStr, nil
}
