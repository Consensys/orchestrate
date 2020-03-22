// +build unit

package dynamic

import (
	"testing"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	mock := &Mock{}
	m := &Middleware{
		Mock: mock,
	}
	assert.Equal(t, "Mock", m.Type())
	field, err := m.Field()
	assert.NoError(t, err)
	assert.Equal(t, mock, field)

	middleware := &traefikdynamic.Middleware{
		StripPrefix: &traefikdynamic.StripPrefix{},
	}
	m = &Middleware{
		Middleware: middleware,
	}
	assert.Equal(t, "StripPrefix", m.Type())
	field, err = m.Field()
	assert.NoError(t, err)
	assert.Equal(t, middleware, field)
}
