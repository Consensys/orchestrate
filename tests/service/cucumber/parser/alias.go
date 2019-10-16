package parser

import (
	"fmt"
	"sync"
)

// AliasRegistry allows to store aliases
// It is safe for concurrent usage
type AliasRegistry struct {
	// aliases[key]value
	mux      *sync.RWMutex
	registry map[string]string
}

func NewAliasRegistry() *AliasRegistry {
	return &AliasRegistry{
		mux:      &sync.RWMutex{},
		registry: make(map[string]string),
	}
}

func (r *AliasRegistry) keyOf(namespace, aka string) string {
	return fmt.Sprintf("%v-%v", namespace, aka)
}

// Get alias from given name space
func (r *AliasRegistry) Get(namespace, aka string) (string, bool) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	if alias, ok := r.registry[r.keyOf(namespace, aka)]; ok {
		return alias, true
	}

	return "", false
}

// Set an alias for a given namespace
func (r *AliasRegistry) Set(namespace, aka, value string) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.registry[r.keyOf(namespace, aka)] = value
}
