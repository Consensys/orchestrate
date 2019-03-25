package testutils

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore"
)

// SecretStoreTestSuite is a test suit for TraceStore
type SecretStoreTestSuite struct {
	suite.Suite
	Store secretstore.SecretStore
}

// TestSecretStore test SecretStore
func (suite *SecretStoreTestSuite) TestSecretStore() {
	err := suite.Store.Store("test-key", "test-value")
	assert.Nil(suite.T(), err, "Secret should have been stored")

	value, ok, err := suite.Store.Load("test-key")
	assert.Nil(suite.T(), err, "Secret should have been loaded")
	assert.True(suite.T(), ok, "Secret should be availaible")
	assert.Equal(suite.T(), "test-value", value, "Secret value should be correct")

	value, ok, err = suite.Store.Load("test-unknown-key")
	assert.Nil(suite.T(), err, "Secret should have been loaded")
	assert.False(suite.T(), ok, "Secret should not be availaible")

	list, err := suite.Store.List()
	assert.Nil(suite.T(), err, "List should be retrieved properly")
	assert.Equal(suite.T(), []string{"test-key"}, list, "Secret list should be correct")

	err = suite.Store.Delete("test-key")
	assert.Nil(suite.T(), err, "Delete should have appened properly")

	list, err = suite.Store.List()
	assert.Nil(suite.T(), err, "List should be retrieved properly")
	assert.Equal(suite.T(), []string{}, list, "Secret list should be correct")

}
