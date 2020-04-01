// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/postgres/migrations"
)

type artifactTestSuite struct {
	suite.Suite
	dataagent           store.ArtifactDataAgent
	tagDataAgent        store.TagDataAgent
	repositoryDataAgent store.RepositoryDataAgent
	pg                  *pgTestUtils.PGTestHelper
}

func TestPGArtifact(t *testing.T) {
	s := new(artifactTestSuite)
	suite.Run(t, s)
}

func (s *artifactTestSuite) SetupSuite() {
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *artifactTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())
	s.dataagent = NewPGArtifact(s.pg.DB)
	s.tagDataAgent = NewPGTag(s.pg.DB)
	s.repositoryDataAgent = NewPGRepository(s.pg.DB)
}

func (s *artifactTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *artifactTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *artifactTestSuite) TestPGArtifact_SelectOrInsertTx() {
	s.T().Run("should insert model successfully", func(t *testing.T) {
		artifact := &models.ArtifactModel{
			Abi:              "ABI",
			Bytecode:         "Bytecode",
			DeployedBytecode: "DeployedBytecode",
			Codehash:         "codeHash",
		}
		err := s.dataagent.SelectOrInsert(context.Background(), artifact)

		assert.Nil(t, err)
		assert.Equal(t, artifact.ID, 1)
	})

	s.T().Run("should insert model successfully in TX", func(t *testing.T) {
		tx, _ := s.pg.DB.Begin()
		ctx := postgres.WithTx(context.Background(), tx)

		artifact := &models.ArtifactModel{
			Abi:              "ABI",
			Bytecode:         "Bytecode",
			DeployedBytecode: "DeployedBytecode",
			Codehash:         "codeHash",
		}
		err := s.dataagent.SelectOrInsert(ctx, artifact)
		_ = tx.Commit()

		assert.Nil(t, err)
		assert.Equal(t, artifact.ID, 1)
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		artifact := &models.ArtifactModel{
			Abi:              "ABI",
			Bytecode:         "Bytecode",
			DeployedBytecode: "DeployedBytecode",
			Codehash:         "codeHash",
		}
		err := s.dataagent.SelectOrInsert(context.Background(), artifact)

		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *artifactTestSuite) TestPGArtifact_FindOneByNameAndTag() {
	s.T().Run("should return NotFoundError if none is found", func(t *testing.T) {
		_, err := s.dataagent.FindOneByNameAndTag(context.Background(), "name", "tag")
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should find successfully", func(t *testing.T) {
		_ = s.insertArtifact("myContract", "tag")

		artifact, err := s.dataagent.FindOneByNameAndTag(context.Background(), "myContract", "tag")

		assert.Nil(t, err)
		assert.Equal(t, 1, artifact.ID)
	})

	s.T().Run("should return PostgresConnectionError if select fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)
		_, err := s.dataagent.FindOneByNameAndTag(context.Background(), "name", "tag")

		assert.True(t, errors.IsPostgresConnectionError(err))

		s.pg.InitTestDB(t)
	})
}

func (s *artifactTestSuite) insertArtifact(name, tag string) error {
	repo := &models.RepositoryModel{
		Name: name,
	}
	_ = s.repositoryDataAgent.SelectOrInsert(context.Background(), repo)

	artifact := &models.ArtifactModel{
		Abi:              "ABI",
		Bytecode:         "Bytecode",
		DeployedBytecode: "DeployedBytecode",
		Codehash:         "codeHash",
	}
	_ = s.dataagent.SelectOrInsert(context.Background(), artifact)

	tagModel := &models.TagModel{
		Name:         tag,
		ArtifactID:   artifact.ID,
		RepositoryID: repo.ID,
	}

	return s.tagDataAgent.Insert(context.Background(), tagModel)
}
