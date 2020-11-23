// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/postgres/migrations"
)

type FaucetModelsTestSuite struct {
	pg *pgTestUtils.PGTestHelper
	FaucetTestSuite
}

func TestModelsFaucet(t *testing.T) {
	s := new(FaucetModelsTestSuite)
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	suite.Run(t, s)
}

func (s *FaucetModelsTestSuite) SetupSuite() {
	s.pg.InitTestDB(s.T())
	s.FaucetAgent = NewPGFaucetAgent(s.pg.DB)
}

func (s *FaucetModelsTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
}

func (s *FaucetModelsTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *FaucetModelsTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

// FaucetTestSuite is a test suite for FaucetRegistry
type FaucetTestSuite struct {
	suite.Suite
	FaucetAgent *PGFaucetAgent
}

const (
	faucetName1 = "testFaucet1"
	faucetName2 = "testFaucet2"
	faucetName3 = "testFaucet3"
)

var tenantID1Faucets = map[string]*models.Faucet{
	faucetName1: {
		Name:            faucetName1,
		TenantID:        tenantID1,
		ChainRule:       "public/",
		CreditorAccount: "0x8dd688660ec0BaBD0B8a2f2DE3232645F73cC5eb",
		MaxBalance:      "1000",
		Amount:          "1000",
		Cooldown:        "1s",
	},
	faucetName2: {
		Name:            faucetName2,
		TenantID:        tenantID1,
		ChainRule:       "public/",
		CreditorAccount: "0x8dd688660ec0BaBD0B8a2f2DE3232645F73cC5eb",
		MaxBalance:      "1000",
		Amount:          "1000",
		Cooldown:        "1s",
	},
}
var tenantID2Faucets = map[string]*models.Faucet{
	faucetName1: {
		Name:            faucetName1,
		TenantID:        tenantID2,
		ChainRule:       "public/",
		CreditorAccount: "0x8dd688660ec0BaBD0B8a2f2DE3232645F73cC5eb",
		MaxBalance:      "1000",
		Amount:          "1000",
		Cooldown:        "1s",
	},
	faucetName2: {
		Name:            faucetName2,
		TenantID:        tenantID2,
		ChainRule:       "public/",
		CreditorAccount: "0x8dd688660ec0BaBD0B8a2f2DE3232645F73cC5eb",
		MaxBalance:      "1000",
		Amount:          "1000",
		Cooldown:        "1s",
	},
	faucetName3: {
		Name:            faucetName3,
		TenantID:        tenantID2,
		ChainRule:       "public/",
		CreditorAccount: "0x8dd688660ec0BaBD0B8a2f2DE3232645F73cC5eb",
		MaxBalance:      "1000",
		Amount:          "1000",
		Cooldown:        "1s",
	},
}

var FaucetsSample = map[string]map[string]*models.Faucet{
	tenantID1: tenantID1Faucets,
	tenantID2: tenantID2Faucets,
}

func CompareFaucets(t *testing.T, faucet1, faucet2 *models.Faucet) {
	assert.Equal(t, faucet1.Name, faucet2.Name, "Should get the same faucet name")
	assert.Equal(t, faucet1.TenantID, faucet2.TenantID, "Should get the same faucet tenantID")
	assert.Equal(t, faucet1.ChainRule, faucet2.ChainRule, "Should get the same faucet ChainRule")
	assert.Equal(t, faucet1.CreditorAccount, faucet2.CreditorAccount, "Should get the same faucet CreditorAccount")
	assert.Equal(t, faucet1.MaxBalance, faucet2.MaxBalance, "Should get the same faucet MaxBalance")
	assert.Equal(t, faucet1.Amount, faucet2.Amount, "Should get the same faucet Amount")
	assert.Equal(t, faucet1.Cooldown, faucet2.Cooldown, "Should get the same faucet Cooldown")
}

func (s *FaucetTestSuite) TestRegisterFaucet() {
	err := s.FaucetAgent.RegisterFaucet(context.Background(), FaucetsSample[tenantID1][faucetName1])
	assert.NoError(s.T(), err, "Should register faucet properly")

	err = s.FaucetAgent.RegisterFaucet(context.Background(), FaucetsSample[tenantID1][faucetName1])
	assert.Error(s.T(), err, "Should get an error violating the 'unique' constraint")
}

func (s *FaucetTestSuite) TestRegisterFaucetWithError() {
	faucetError := &models.Faucet{
		Name:       "faucetName1",
		TenantID:   "tenantID1",
		ChainRule:  "public/",
		MaxBalance: "1000",
		Amount:     "1000",
		Cooldown:   "1s",
	}
	err := s.FaucetAgent.RegisterFaucet(context.Background(), faucetError)
	assert.Error(s.T(), err, "Should get an error when a field is missing")
}

func (s *FaucetTestSuite) TestRegisterFaucets() {
	var err error
	for _, faucets := range FaucetsSample {
		for _, faucet := range faucets {
			err = s.FaucetAgent.RegisterFaucet(context.Background(), faucet)
			assert.NoError(s.T(), err, "should not error on registration")
		}
	}
}

func (s *FaucetTestSuite) TestGetFaucets() {
	s.TestRegisterFaucets()

	faucets, err := s.FaucetAgent.GetFaucets(context.Background(), nil, nil)
	assert.NoError(s.T(), err, "Should get faucets without errors")
	assert.Len(s.T(), faucets, len(tenantID1Faucets)+len(tenantID2Faucets), "Should get the same number of faucets")

	for _, faucet := range faucets {
		CompareFaucets(s.T(), faucet, FaucetsSample[faucet.TenantID][faucet.Name])
	}
}

func (s *FaucetTestSuite) TestGetFaucetByUUID() {
	s.TestRegisterFaucets()

	faucetUUID := FaucetsSample[tenantID1][faucetName1].UUID

	faucet, err := s.FaucetAgent.GetFaucet(context.Background(), faucetUUID, nil)
	assert.NoError(s.T(), err, "Should get faucet without errors")

	CompareFaucets(s.T(), faucet, FaucetsSample[tenantID1][faucetName1])
}

func (s *FaucetTestSuite) TestGetFaucetByUUIDByTenant() {
	s.TestRegisterFaucets()

	faucetUUID := FaucetsSample[tenantID1][faucetName1].UUID

	faucet, err := s.FaucetAgent.GetFaucet(context.Background(), faucetUUID, []string{tenantID1})
	assert.NoError(s.T(), err, "Should get faucet without errors")

	assert.Equal(s.T(), tenantID1, faucet.TenantID)
}

func (s *FaucetTestSuite) TestUpdateFaucetByUUID() {
	s.TestRegisterFaucets()

	testFaucet := FaucetsSample[tenantID1][faucetName2]
	testFaucet.ChainRule = "private"
	err := s.FaucetAgent.UpdateFaucet(context.Background(), testFaucet.UUID, nil, testFaucet)
	assert.NoError(s.T(), err, "Should update faucet without errors")

	faucet, _ := s.FaucetAgent.GetFaucet(context.Background(), testFaucet.UUID, nil)
	CompareFaucets(s.T(), faucet, testFaucet)
}

func (s *FaucetTestSuite) TestErrorNotFoundUpdateFaucetByUUID() {
	s.TestRegisterFaucets()

	testFaucet := &models.Faucet{
		UUID:      "0d60a85e-0b90-4482-a14c-108aea2557aa",
		ChainRule: "private",
	}
	err := s.FaucetAgent.UpdateFaucet(context.Background(), testFaucet.UUID, nil, testFaucet)
	assert.Error(s.T(), err, "Should update faucet with errors")
}

func (s *FaucetTestSuite) TestDeleteFaucetByUUID() {
	s.TestRegisterFaucets()

	faucetUUID := FaucetsSample[tenantID1][faucetName1].UUID

	err := s.FaucetAgent.DeleteFaucet(context.Background(), faucetUUID, nil)
	assert.NoError(s.T(), err, "Should delete faucet without errors")
}

func (s *FaucetTestSuite) TestDeleteFaucetByUUIDByTenant() {
	s.TestRegisterFaucets()

	faucetUUID := FaucetsSample[tenantID1][faucetName1].UUID

	err := s.FaucetAgent.DeleteFaucet(context.Background(), faucetUUID, []string{tenantID1})
	assert.NoError(s.T(), err, "Should delete faucet without errors")
}

func (s *FaucetTestSuite) TestErrorNotFoundDeleteFaucetByUUIDAndTenant() {
	s.TestRegisterFaucets()

	// tenantID2 in the context but we try to delete the faucetUUID of tenantID1
	faucetUUID := FaucetsSample[tenantID1][faucetName1].UUID

	err := s.FaucetAgent.DeleteFaucet(context.Background(), faucetUUID, []string{tenantID2})
	assert.Error(s.T(), err, "Should delete faucet with errors")
}

func (s *FaucetTestSuite) TestErrorNotFoundDeleteFaucetByUUID() {
	s.TestRegisterFaucets()

	err := s.FaucetAgent.DeleteFaucet(context.Background(), "0d60a85e-0b90-4482-a14c-108aea2557aa", nil)
	assert.Error(s.T(), err, "Should delete faucet with errors")
}
