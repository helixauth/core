package secrets

import (
	"fmt"
	"strings"
)

// Get retrieves a secret for the provided key
func (m *manager) Get(key string) (interface{}, error) {
	keys := strings.Split(key, ".")
	secrets := m.secrets
	for i, k := range keys {
		if i == len(keys)-1 {
			return secrets[k], nil
		}
		ok := true
		if secrets, ok = secrets[k].(map[interface{}]interface{}); !ok {
			return nil, fmt.Errorf("No secret for key '%v", key)
		}
	}
	return "", fmt.Errorf("Invalid key '%v", key)
}
