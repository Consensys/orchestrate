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

const (
	abi = "contractABI"
)

type artifactTestSuite struct {
	suite.Suite
	agents *PGAgents
	pg     *pgTestUtils.PGTestHelper
}

func TestPGArtifact(t *testing.T) {
	s := new(artifactTestSuite)
	suite.Run(t, s)
}

func (s *artifactTestSuite) SetupSuite() {
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *artifactTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.agents = New(s.pg.DB)
}

func (s *artifactTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *artifactTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *artifactTestSuite) TestPGArtifact_Insert() {
	ctx := context.Background()

	s.T().Run("should insert model successfully", func(t *testing.T) {
		artifact := &models.ArtifactModel{
			ABI:              abi,
			Bytecode:         "0x123",
			DeployedBytecode: "0x123",
			Codehash:         codeHash,
		}
		err := s.agents.Artifact().Insert(ctx, artifact)

		assert.NoError(t, err)
		assert.NotEmpty(t, artifact.ID, 1)
	})

	s.T().Run("should insert model successfully in TX", func(t *testing.T) {
		tx, _ := s.pg.DB.Begin()
		ctx2 := postgres.WithTx(ctx, tx)

		artifact := &models.ArtifactModel{
			ABI:              "ABI2",
			Bytecode:         "0x321",
			DeployedBytecode: "0x321",
			Codehash:         codeHash,
		}
		err := s.agents.Artifact().Insert(ctx2, artifact)
		_ = tx.Commit()

		assert.NoError(t, err)
		assert.NotEmpty(t, artifact.ID)
	})
	
	s.T().Run("should fail to insert model duplicated model", func(t *testing.T) {

		artifact := &models.ArtifactModel{
			ABI:              abi,
			Bytecode:         "0x123",
			DeployedBytecode: "0x123",
			Codehash:         codeHash,
		}
		err := s.agents.Artifact().Insert(ctx, artifact)

		assert.Error(t, err)
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		artifact := &models.ArtifactModel{
			ABI:              "ABI",
			Bytecode:         "Bytecode",
			DeployedBytecode: "DeployedBytecode",
			Codehash:         "codeHash",
		}
		err := s.agents.Artifact().Insert(ctx, artifact)

		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *artifactTestSuite) TestPGArtifact_FindOneByABIAndCodeHash() {
	ctx := context.Background()

	s.T().Run("should return NotFoundError if none is found", func(t *testing.T) {
		_, err := s.agents.Artifact().FindOneByABIAndCodeHash(ctx, abi, codeHash)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should find successfully", func(t *testing.T) {
		_ = s.insertArtifacts(ctx, "myContract", "tag")

		artifact, err := s.agents.Artifact().FindOneByABIAndCodeHash(ctx, abi, codeHash)

		assert.NoError(t, err)
		assert.NotEmpty(t, artifact.ID)
	})
	
	s.T().Run("should return PostgresConnectionError if select fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)
		_, err := s.agents.Artifact().FindOneByABIAndCodeHash(ctx, abi, codeHash)

		assert.True(t, errors.IsPostgresConnectionError(err))

		s.pg.InitTestDB(t)
	})
}

func (s *artifactTestSuite) TestPGArtifact_FindOneByNameAndTag() {
	ctx := context.Background()

	s.T().Run("should return NotFoundError if none is found", func(t *testing.T) {
		_, err := s.agents.Artifact().FindOneByNameAndTag(ctx, "name", "tag")
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should find successfully", func(t *testing.T) {
		_ = s.insertArtifacts(ctx, "myContract", "tag")

		artifact, err := s.agents.Artifact().FindOneByNameAndTag(ctx, "myContract", "tag")

		assert.NoError(t, err)
		assert.Equal(t, 1, artifact.ID)
	})

	s.T().Run("should return PostgresConnectionError if select fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)
		_, err := s.agents.Artifact().FindOneByNameAndTag(ctx, "name", "tag")

		assert.True(t, errors.IsPostgresConnectionError(err))

		s.pg.InitTestDB(t)
	})
}

func (s *artifactTestSuite) insertArtifacts(ctx context.Context, name, tag string) error {
	repo := &models.RepositoryModel{
		Name: name,
	}
	_ = s.agents.Repository().Insert(ctx, repo)

	artifact := &models.ArtifactModel{
		ABI:              abi,
		Bytecode:         "0x234",
		DeployedBytecode: "0x234",
		Codehash:         codeHash,
	}
	_ = s.agents.Artifact().Insert(ctx, artifact)

	tagModel := &models.TagModel{
		Name:         tag,
		ArtifactID:   artifact.ID,
		RepositoryID: repo.ID,
	}

	return s.agents.Tag().Insert(ctx, tagModel)
}
