package testutils

import (
	"context"
	"fmt"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

var httpRouter1 = &dynamic.Router{
	Service: "testService",
}
var JSONHttpRouter1, _ = json.Marshal(httpRouter1)

var httpRouter2 = &dynamic.Router{
	Rule: "testRule",
}
var JSONHttpRouter2, _ = json.Marshal(httpRouter2)

var httpMiddleware1 = &dynamic.Middleware{
	AddPrefix: &dynamic.AddPrefix{Prefix: "testPrefix"},
}
var JSONHttpMiddleware1, _ = json.Marshal(httpMiddleware1)

var configs = []types.Config{
	{
		Name:       "testRouter1",
		ConfigType: types.HTTPROUTER,
		Config:     JSONHttpRouter1,
	},
	{
		Name:       "testRouter2",
		ConfigType: types.HTTPROUTER,
		Config:     JSONHttpRouter2,
	},
	{
		Name:       "testMiddleware1",
		ConfigType: types.HTTPMIDDLEWARE,
		Config:     JSONHttpMiddleware1,
	},
}

var configsWithDifferentTenants = []types.Config{
	{
		Name:       "testRouter1",
		TenantID:   "tenant1",
		ConfigType: types.HTTPROUTER,
		Config:     JSONHttpRouter1,
	},
	{
		Name:       "testRouter2",
		TenantID:   "tenant2",
		ConfigType: types.HTTPROUTER,
		Config:     JSONHttpRouter2,
	},
	{
		Name:       "testMiddleware1",
		TenantID:   "tenant2",
		ConfigType: types.HTTPMIDDLEWARE,
		Config:     JSONHttpMiddleware1,
	},
	{
		Name:       "testMiddleware1",
		TenantID:   "tenant1",
		ConfigType: types.HTTPMIDDLEWARE,
		Config:     JSONHttpMiddleware1,
	},
}

// EnvelopeStoreTestSuite is a test suit for EnvelopeStore
type ChainRegistryTestSuite struct {
	suite.Suite
	Store types.ChainRegistryStore
}

func (s *ChainRegistryTestSuite) TestRegisterConfig() {
	config := &types.Config{
		Name:       "test",
		ConfigType: types.HTTPROUTER,
		Config:     JSONHttpRouter1,
	}

	err := s.Store.RegisterConfig(context.Background(), config)
	assert.NoError(s.T(), err, "Should register config properly")

	err = s.Store.RegisterConfig(context.Background(), config)
	assert.Error(s.T(), err, "Should get an error violating the 'unique' constrain")

	errorConfig := &types.Config{
		Name:       "test",
		ConfigType: types.HTTPMIDDLEWARE,
		Config:     JSONHttpRouter1,
	}
	err = s.Store.RegisterConfig(context.Background(), errorConfig)
	assert.Error(s.T(), err, "Should get an error when config does not corresponds to the config type")
}

func (s *ChainRegistryTestSuite) TestRegisterConfigs() {
	err := s.Store.RegisterConfigs(context.Background(), &configs)
	assert.NoError(s.T(), err, "Should register configs properly")
}

func (s *ChainRegistryTestSuite) TestRegisterConfigsWithDifferentTenant() {
	err := s.Store.RegisterConfigs(context.Background(), &configsWithDifferentTenants)
	assert.NoError(s.T(), err, "Should register configs properly")
}

func (s *ChainRegistryTestSuite) TestGetConfigByID() {
	s.TestRegisterConfigs()
	config := &types.Config{
		ID: 2,
	}
	err := s.Store.GetConfigByID(context.Background(), config)
	assert.NoError(s.T(), err, "Should get config properly")
	assert.Equal(s.T(), configs[config.ID-1].Name, config.Name, "GetConfigById should retrieve the correct name")
	assert.Equal(s.T(), configs[config.ID-1].ConfigType, config.ConfigType, "GetConfigById should retrieve the correct configType")

	configError := &types.Config{
		ID: -1,
	}
	err = s.Store.GetConfigByID(context.Background(), configError)
	assert.Error(s.T(), err, "Should get an error when ID is wrong")
}

func (s *ChainRegistryTestSuite) TestGetConfigByTenantID() {
	s.TestRegisterConfigs()
	config := &types.Config{
		TenantID: "default",
	}
	configs, err := s.Store.GetConfigByTenantID(context.Background(), config)
	assert.NoError(s.T(), err, "Should get configs properly")
	assert.Len(s.T(), configs, len(configs), fmt.Sprintf("Should get %d configs", len(configs)))
}

func (s *ChainRegistryTestSuite) TestUpdateConfigByID() {
	s.TestRegisterConfigs()

	JSONHttpMiddleware2, _ := json.Marshal(&dynamic.Middleware{
		AddPrefix: &dynamic.AddPrefix{Prefix: "newTestPrefix"},
	})

	config := &types.Config{
		ID:         3,
		Name:       "testMiddleware2",
		ConfigType: types.HTTPMIDDLEWARE,
		Config:     JSONHttpMiddleware2,
	}
	err := s.Store.UpdateConfigByID(context.Background(), config)
	assert.NoError(s.T(), err, "Should update config properly")

	configReq := &types.Config{ID: 3}
	err = s.Store.GetConfigByID(context.Background(), configReq)
	assert.NoError(s.T(), err, "Should get config properly")
	assert.Equal(s.T(), configReq.Name, config.Name, "GetConfigById should retrieve the correct name")
	assert.Equal(s.T(), configReq.ConfigType, config.ConfigType, "GetConfigById should retrieve the correct name")
}

func (s *ChainRegistryTestSuite) TestDeregisterConfigByID() {
	s.TestRegisterConfigs()

	err := s.Store.DeregisterConfigByID(context.Background(), &types.Config{ID: 1})
	assert.NoError(s.T(), err, "Should get config properly")

	configReq := &types.Config{
		TenantID: "default",
	}
	configsReq, err := s.Store.GetConfigByTenantID(context.Background(), configReq)
	assert.NoError(s.T(), err, "Should get configs properly")
	assert.Len(s.T(), configsReq, len(configs)-1, "")
}

func (s *ChainRegistryTestSuite) TestDeregisterConfigsByIds() {
	s.TestRegisterConfigs()

	var configsToDelete = []types.Config{{ID: 1}, {ID: 2}}

	err := s.Store.DeregisterConfigsByIds(context.Background(), &configsToDelete)
	assert.NoError(s.T(), err, "Should get config properly")

	configReq := &types.Config{
		TenantID: "default",
	}
	configsReq, err := s.Store.GetConfigByTenantID(context.Background(), configReq)
	assert.NoError(s.T(), err, "Should get configs properly")
	assert.Len(s.T(), configsReq, len(configs)-len(configsToDelete), "")
}

func (s *ChainRegistryTestSuite) TestDeregisterConfigByTenantID() {
	s.TestRegisterConfigsWithDifferentTenant()
	tenantIDToDelete := "tenant1"

	configsReq, err := s.Store.GetConfigByTenantID(context.Background(), &types.Config{TenantID: tenantIDToDelete})
	assert.NoError(s.T(), err, "Should get configs properly")
	assert.Len(s.T(), configsReq, 2, "Should have config")

	err = s.Store.DeregisterConfigByTenantID(context.Background(), &types.Config{TenantID: tenantIDToDelete})
	assert.NoError(s.T(), err, "Should delete config properly")

	configsReq2, err := s.Store.GetConfigByTenantID(context.Background(), &types.Config{TenantID: tenantIDToDelete})
	assert.NoError(s.T(), err, "Should get configs properly")
	assert.Len(s.T(), configsReq2, 0, "Should not have config anymore")
}
