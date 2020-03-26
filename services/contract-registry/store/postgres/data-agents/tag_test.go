// +build !race

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

type tagTestSuite struct {
	suite.Suite
	dataagent  store.TagDataAgent
	repoDA     store.RepositoryDataAgent
	artifactDA store.ArtifactDataAgent
	pg         *pgTestUtils.PGTestHelper
}

func TestPGTag(t *testing.T) {
	s := new(tagTestSuite)
	suite.Run(t, s)
}

func (s *tagTestSuite) SetupSuite() {
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *tagTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())

	s.repoDA = NewPGRepository(s.pg.DB)
	s.artifactDA = NewPGArtifact(s.pg.DB)

	s.dataagent = NewPGTag(s.pg.DB)
}

func (s *tagTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *tagTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *tagTestSuite) TestPGTag_Insert() {
	s.T().Run("should insert model successfully", func(t *testing.T) {
		tag, err := s.insertTag(context.Background(), "contract", "tag")

		assert.Nil(t, err)
		assert.Equal(t, 1, tag.ID)
	})

	s.T().Run("should insert model successfully in TX", func(t *testing.T) {
		tx, _ := s.pg.DB.Begin()
		ctx := postgres.WithTx(context.Background(), tx)

		tag, err := s.insertTag(ctx, "contract", "tag")
		_ = tx.Commit()

		assert.Nil(t, err)
		assert.NotEmpty(t, tag.ID)
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		tag := &models.TagModel{
			Name:         "tagName",
			RepositoryID: 1,
			ArtifactID:   1,
		}
		err := s.dataagent.Insert(context.Background(), tag)

		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *tagTestSuite) TestPGRepository_FindAllByName() {
	contractName := "myContract"

	s.T().Run("should return NotFoundError if none is found", func(t *testing.T) {
		_, err := s.dataagent.FindAllByName(context.Background(), contractName)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should find all successfully", func(t *testing.T) {
		_, _ = s.insertTag(context.Background(), contractName, "tag")

		tags, err := s.dataagent.FindAllByName(context.Background(), contractName)

		assert.Equal(t, 1, len(tags))
		assert.Equal(t, "tag", tags[0])
		assert.Nil(t, err)
	})

	s.T().Run("should return PostgresConnectionError if select fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)
		_, err := s.dataagent.FindAllByName(context.Background(), contractName)

		assert.True(t, errors.IsPostgresConnectionError(err))

		s.pg.InitTestDB(t)
	})
}

func (s *tagTestSuite) insertTag(ctx context.Context, contractName, tagName string) (*models.TagModel, error) {
	repo := &models.RepositoryModel{
		Name: contractName,
	}
	_ = s.repoDA.SelectOrInsert(ctx, repo)

	artifact := &models.ArtifactModel{
		Abi:              "ABI",
		Bytecode:         "Bytecode",
		DeployedBytecode: "DeployedBytecode",
		Codehash:         "codeHash",
	}
	_ = s.artifactDA.SelectOrInsert(ctx, artifact)

	tag := &models.TagModel{
		Name:         tagName,
		RepositoryID: repo.ID,
		ArtifactID:   artifact.ID,
	}

	err := s.dataagent.Insert(ctx, tag)

	return tag, err
}
