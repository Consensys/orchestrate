package alias

import (
	"reflect"
	"strings"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/utils"
)

// Registry allows to store aliases
// It is safe for concurrent usage
type Registry struct {
	// aliases[key]value
	mux      *sync.RWMutex
	registry map[string]interface{}
	logger   log.Logger
}

func NewAliasRegistry() *Registry {
	logger := log.WithoutContext()
	return &Registry{
		mux:      &sync.RWMutex{},
		registry: make(map[string]interface{}),
		logger:   logger,
	}
}

// Get alias from given name space
func (r *Registry) Get(aka ...string) (interface{}, bool) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	v, err := utils.GetField(strings.Join(aka, "."), reflect.ValueOf(r.registry))
	if err != nil || v.Kind() == reflect.Invalid {
		return nil, false
	}
	return v.Interface(), true
}

// Set an alias for a given namespace
func (r *Registry) Set(value interface{}, aka ...string) bool {
	r.mux.Lock()
	defer r.mux.Unlock()

	ok := setMap(strings.Join(aka, "."), value, r.registry)
	if ok {
		r.logger.
			WithField("aka", aka).
			WithField("value", value).
			Debugf("Registry value set")
	}

	return ok
}

func setMap(path string, value interface{}, m map[string]interface{}) bool {
	key := strings.Split(path, ".")
	if len(key) == 1 {
		m[path] = value
		return true
	}
	field, ok := m[key[0]]
	if !ok {
		newNestedMap := make(map[string]interface{})
		m[key[0]] = newNestedMap
		return setMap(strings.Join(key[1:], "."), value, newNestedMap)
	}

	if nestedMap, ok := field.(map[string]interface{}); ok {
		return setMap(strings.Join(key[1:], "."), value, nestedMap)
	}
	return false
}
