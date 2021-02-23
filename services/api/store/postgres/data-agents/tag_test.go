// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/postgres/migrations"
)

type tagTestSuite struct {
	suite.Suite
	agents *PGAgents
	pg     *pgTestUtils.PGTestHelper
}

func TestPGTag(t *testing.T) {
	s := new(tagTestSuite)
	suite.Run(t, s)
}

func (s *tagTestSuite) SetupSuite() {
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *tagTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.agents = New(s.pg.DB)
}

func (s *tagTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *tagTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *tagTestSuite) TestPGTag_Insert() {
	ctx := context.Background()
	s.T().Run("should insert model successfully", func(t *testing.T) {
		tag, err := s.insertTag(ctx, "contract", "tag")

		assert.NoError(t, err)
		assert.Equal(t, 1, tag.ID)
	})

	s.T().Run("should insert model successfully in TX", func(t *testing.T) {
		tx, _ := s.pg.DB.Begin()
		ctx2 := postgres.WithTx(ctx, tx)

		tag, err := s.insertTag(ctx2, "contract", "tag")
		_ = tx.Commit()

		assert.NoError(t, err)
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
		err := s.agents.Tag().Insert(ctx, tag)

		assert.True(t, errors.IsInternalError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *tagTestSuite) TestPGTag_FindAllByName() {
	ctx := context.Background()
	contractName := "myContract"

	s.T().Run("should return NotFoundError if none is found", func(t *testing.T) {
		result, err := s.agents.Tag().FindAllByName(ctx, contractName)
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	s.T().Run("should find all successfully", func(t *testing.T) {
		_, _ = s.insertTag(ctx, contractName, "tag")

		tags, err := s.agents.Tag().FindAllByName(ctx, contractName)

		assert.Equal(t, 1, len(tags))
		assert.Equal(t, "tag", tags[0])
		assert.NoError(t, err)
	})

	s.T().Run("should return PostgresConnectionError if select fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)
		_, err := s.agents.Tag().FindAllByName(ctx, contractName)

		assert.True(t, errors.IsInternalError(err))

		s.pg.InitTestDB(t)
	})
}

func (s *tagTestSuite) insertTag(ctx context.Context, contractName, tagName string) (*models.TagModel, error) {
	repo := &models.RepositoryModel{
		Name: contractName,
	}
	_ = s.agents.Repository().Insert(ctx, repo)

	artifact := &models.ArtifactModel{
		ABI:              abi,
		Bytecode:         "Bytecode",
		DeployedBytecode: "DeployedBytecode",
		Codehash:         codeHash,
	}
	_ = s.agents.Artifact().Insert(ctx, artifact)

	tag := &models.TagModel{
		Name:         tagName,
		RepositoryID: repo.ID,
		ArtifactID:   artifact.ID,
	}

	err := s.agents.Tag().Insert(ctx, tag)

	return tag, err
}
