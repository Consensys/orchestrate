// +build unit

package dynamic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	healthcheck := &HealthCheck{}
	s := &Service{
		HealthCheck: healthcheck,
	}
	assert.Equal(t, "HealthCheck", s.Type())
	field, err := s.Field()
	assert.NoError(t, err)
	assert.Equal(t, healthcheck, field)

	proxy := &ReverseProxy{}
	s = &Service{
		ReverseProxy: proxy,
	}

	assert.Equal(t, "ReverseProxy", s.Type())
	field, err = s.Field()
	assert.NoError(t, err)
	assert.Equal(t, proxy, field)
}
