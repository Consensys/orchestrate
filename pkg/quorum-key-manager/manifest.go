package quorumkeymanager

import (
	"gopkg.in/yaml.v2"
)

// Manifest for a store
type Manifest struct {
	// Kind  of store
	Kind string `json:"kind" validate:"required"`

	// Version
	Version string `json:"version"`

	// Name of the store
	Name string `json:"name" validate:"required"`

	// Tags are user set information about a store
	Tags map[string]string `json:"tags"`

	// Specs specifications of a store
	Specs interface{} `json:"specs" validate:"required"`
}

func (m *Manifest) MarshallToYaml() ([]byte, error) {
	b, err := yaml.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}
