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

// Init initializes key Envelope with EnabledViperKey
func Init(_ context.Context) {
	initOnce.Do(func() {
		if keyBuilder != nil {
			return
		}

		keyBuilder = New(viper.GetBool(EnabledViperKey))
	})
}

// GlobalKeyBuilder returns global Authentication Manager
func GlobalKeyBuilder() *KeyBuilder {
	return keyBuilder
}

// SetKeyBuilder sets global Authentication Manager
func SetKeyBuilder(key *KeyBuilder) {
	keyBuilder = key
}
