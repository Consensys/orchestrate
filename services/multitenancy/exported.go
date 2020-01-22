package multitenancy

import (
	"context"
	"sync"

	"github.com/spf13/viper"
)

var (
	keyBuilder *KeyBuilder
	initOnce   = &sync.Once{}
)

const (
	DefaultTenantIDName = "_"
)

// Init initializes key Builder with EnabledViperKey
func Init(_ context.Context) {
	initOnce.Do(func() {
		if keyBuilder != nil {
			return
		}

		keyBuilder = New(viper.GetBool(EnabledViperKey))
	})
}

// GlobalAuth returns global Authentication Manager
func GlobalKeyBuilder() *KeyBuilder {
	return keyBuilder
}

// SetGlobalAuth sets global Authentication Manager
func SetKeyBuilder(key *KeyBuilder) {
	keyBuilder = key
}
