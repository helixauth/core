package secrets

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// SecretsManager provides an interface for fetching secrets
type SecretsManager interface {
	Get(key string) (interface{}, error)
}

type manager struct {
	secrets map[interface{}]interface{}
}

// New creates a new SecretsManager
func New(filename string) (SecretsManager, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	secrets := map[interface{}]interface{}{}
	if err = yaml.Unmarshal(yamlFile, &secrets); err != nil {
		return nil, err
	}
	return &manager{
		secrets: secrets,
	}, nil
}

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
