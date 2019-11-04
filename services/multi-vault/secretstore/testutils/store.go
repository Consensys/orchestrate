package testutils

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/secretstore/services"
)

// SecretStoreTestSuite is a test suit for TraceStore
type SecretStoreTestSuite struct {
	suite.Suite
	Store services.SecretStore
}

// TestSecretStore test SecretStore
func (s *SecretStoreTestSuite) TestSecretStore() {

	err := s.Store.Store("test-key", "test-value")
	assert.Nilf(s.T(), err, "Secret should have been stored, got %q", err)

	value, ok, err := s.Store.Load("test-key")
	assert.Nilf(s.T(), err, "Secret should have been loaded, got %q", err)
	assert.True(s.T(), ok, "Secret should be available")
	assert.Equal(s.T(), "test-value", value, "Secret value should be correct")

	_, ok, err = s.Store.Load("test-unknown-key")
	assert.Nilf(s.T(), err, "Secret should have been loaded, got %q", err)
	assert.False(s.T(), ok, "Secret should not be available")

	list, err := s.Store.List()
	assert.Nilf(s.T(), err, "List should be retrieved properly, got %q", err)
	assert.Equal(s.T(), []string{"test-key"}, list, "Secret list should be correct")

	err = s.Store.Delete("test-key")
	assert.Nilf(s.T(), err, "Delete should have happened properly, got %q", err)

	list, err = s.Store.List()
	assert.Nilf(s.T(), err, "List should be retrieved properly, got %q", err)
	assert.Equal(s.T(), []string{}, list, "Secret list should be correct")

	value, ok, err = s.Store.Load("test-key")
	assert.Nilf(s.T(), err, "Load should not return an error, got %q", err)
	assert.Falsef(s.T(), ok, "Secret should not have been found")
	assert.Equal(s.T(), "", value, "Secret list should be correct")

}
