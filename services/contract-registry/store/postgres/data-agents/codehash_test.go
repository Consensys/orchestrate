// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/postgres/migrations"
)

type codeHashTestSuite struct {
	suite.Suite
	dataagent store.CodeHashDataAgent
	pg        *pgTestUtils.PGTestHelper
}

func TestPGCodeHash(t *testing.T) {
	s := new(codeHashTestSuite)
	suite.Run(t, s)
}

func (s *codeHashTestSuite) SetupSuite() {
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *codeHashTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())
	s.dataagent = NewPGCodeHash(s.pg.DB)
}

func (s *codeHashTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *codeHashTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *codeHashTestSuite) TestPGArtifact_Insert() {
	s.T().Run("should insert model successfully", func(t *testing.T) {
		codehash := &models.CodehashModel{
			ChainID:  "chainID",
			Address:  "address",
			Codehash: "codeHash",
		}
		err := s.dataagent.Insert(context.Background(), codehash)

		assert.Nil(t, err)
		assert.Equal(t, codehash.ID, 1)
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		codehash := &models.CodehashModel{
			ChainID:  "chainID",
			Address:  "address",
			Codehash: "codeHash",
		}
		err := s.dataagent.Insert(context.Background(), codehash)

		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}
