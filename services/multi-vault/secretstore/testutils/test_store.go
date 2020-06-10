package testutils

import (
	"context"

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

	err := s.Store.Store(context.Background(), "test-key", "test-value")
	assert.NoError(s.T(), err, "Secret should have been stored")

	err = s.Store.Store(context.Background(), "test-key", "test-value")
	assert.NoError(s.T(), err, "Secret should have been stored even if it already exists")

	err = s.Store.Store(context.Background(), "test-key", "test-value-not-the-same")
	assert.Error(s.T(), err, "Store should fail when trying to save twice the same key")

	value, ok, err := s.Store.Load(context.Background(), "test-key")
	assert.NoError(s.T(), err, "Secret should have been loaded")
	assert.True(s.T(), ok, "Secret should be available")
	assert.Equal(s.T(), "test-value", value, "Secret value should be correct")

	_, ok, err = s.Store.Load(context.Background(), "test-unknown-key")
	assert.NoError(s.T(), err, "Secret should have been loaded")
	assert.False(s.T(), ok, "Secret should not be available")

	list, err := s.Store.List()
	assert.NoError(s.T(), err, "List should be retrieved properly")
	assert.Equal(s.T(), []string{"_test-key"}, list, "Secret list should be correct")

	err = s.Store.Delete(context.Background(), "test-key")
	assert.NoError(s.T(), err, "Delete should have happened properly")

	list, err = s.Store.List()
	assert.NoError(s.T(), err, "List should be retrieved properly")
	assert.Equal(s.T(), []string{}, list, "Secret list should be correct")

	value, ok, err = s.Store.Load(context.Background(), "test-key")
	assert.NoError(s.T(), err, "Load should not return an error")
	assert.Falsef(s.T(), ok, "Secret should not have been found")
	assert.Equal(s.T(), "", value, "Secret list should be correct")
}
