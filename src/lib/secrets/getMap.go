package secrets

import "fmt"

// Get retrieves a secret map for the provided key
func (m *manager) GetMap(key string) (map[interface{}]interface{}, error) {
	sec, err := m.Get(key)
	if err != nil {
		return nil, err
	}
	secMap, ok := sec.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("Secret is not available")
	}
	return secMap, nil
}
