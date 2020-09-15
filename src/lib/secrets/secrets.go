package secrets

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Manager provides an interface for fetching secrets
type Manager interface {
	Get(key string) (interface{}, error)
	GetMap(key string) (map[interface{}]interface{}, error)
	GetString(key string) (string, error)
}

type manager struct {
	secrets map[interface{}]interface{}
}

// New creates a new Manager
func New(filename string) (Manager, error) {
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
