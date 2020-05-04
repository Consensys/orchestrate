package parser

import (
	"fmt"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
)

// AliasRegistry allows to store aliases
// It is safe for concurrent usage
type AliasRegistry struct {
	// aliases[key]value
	mux      *sync.RWMutex
	registry map[string]string
	logger   log.Logger
}

func NewAliasRegistry() *AliasRegistry {
	logger := log.WithoutContext()
	return &AliasRegistry{
		mux:      &sync.RWMutex{},
		registry: make(map[string]string),
		logger:   logger,
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

	r.logger.
		WithField("namespace", namespace).
		WithField("aka", aka).
		WithField("value", value).
		Debugf("AliasRegistry value set")

	r.registry[r.keyOf(namespace, aka)] = value
}
