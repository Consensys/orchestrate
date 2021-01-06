// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/postgres/migrations"
)

type repositoryTestSuite struct {
	suite.Suite
	agents *PGAgents
	pg     *pgTestUtils.PGTestHelper
}

func TestPGRepository(t *testing.T) {
	s := new(repositoryTestSuite)
	suite.Run(t, s)
}

func (s *repositoryTestSuite) SetupSuite() {
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *repositoryTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.agents = New(s.pg.DB)
}

func (s *repositoryTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *repositoryTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *repositoryTestSuite) TestPGRepository_Insert() {
	ctx := context.Background()

	s.T().Run("should insert model successfully", func(t *testing.T) {
		repo := &models.RepositoryModel{
			Name: "myRepository",
		}
		err := s.agents.Repository().Insert(ctx, repo)

		assert.NoError(t, err)
		assert.Equal(t, 1, repo.ID)
	})

	s.T().Run("should insert model successfully in TX", func(t *testing.T) {
		tx, _ := s.pg.DB.Begin()
		ctx2 := postgres.WithTx(ctx, tx)

		repo := &models.RepositoryModel{
			Name: "myRepository2",
		}
		err := s.agents.Repository().Insert(ctx2, repo)
		_ = tx.Commit()

		assert.NoError(t, err)
		assert.NotEmpty(t, repo.ID)
	})
	
	s.T().Run("should fail to insert model duplicated repo", func(t *testing.T) {
		repo := &models.RepositoryModel{
			Name: "myRepository",
		}
		err := s.agents.Repository().Insert(ctx, repo)

		assert.Error(t, err)
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		repo := &models.RepositoryModel{
			Name: "myRepository",
		}
		err := s.agents.Repository().Insert(ctx, repo)

		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *repositoryTestSuite) TestPGRepository_LockOneByName() {
	ctx := context.Background()
	s.insertRepo(ctx, 1)

	s.T().Run("should lock successfully", func(t *testing.T) {
		dbtx0, err := s.pg.DB.Begin()
		assert.NoError(t, err)
		dbtx1, err := s.pg.DB.Begin()
		assert.NoError(t, err)
		newPGRepo0 := NewPGRepository(dbtx0)
		newPGRepo1 := NewPGRepository(dbtx1)

		waitChannel := make(chan string)
		_, err = newPGRepo1.FindOneAndLock(ctx, "myRepository_0")
		assert.NoError(t, err)
		go func() {
			time.Sleep(2000 * time.Millisecond)
			waitChannel <- "repo1"

			err = dbtx1.Commit()
			assert.NoError(t, err)
		}()

		go func() {
			_, err = newPGRepo0.FindOneAndLock(ctx, "myRepository_0")
			assert.NoError(t, err)

			err = dbtx0.Commit()
			assert.NoError(t, err)

			waitChannel <- "repo0"
		}()

		firstRepo := <-waitChannel
		assert.Equal(t, firstRepo, "repo1")

		secondRepo := <-waitChannel
		assert.Equal(t, secondRepo, "repo0")
	})
}

func (s *repositoryTestSuite) TestPGRepository_FindOneByName() {
	ctx := context.Background()

	s.T().Run("should return NotFoundError if none is found", func(t *testing.T) {
		_, err := s.agents.Repository().FindOne(ctx, "unknown")
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should find successfully", func(t *testing.T) {
		s.insertRepo(ctx, 1)

		artifact, err := s.agents.Repository().FindOne(ctx, "myRepository_0")

		assert.NoError(t, err)
		assert.NotEmpty(t, artifact.ID)
	})
	
	s.T().Run("should return PostgresConnectionError if select fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)
		_, err := s.agents.Repository().FindOne(ctx,  "respository")

		assert.True(t, errors.IsPostgresConnectionError(err))

		s.pg.InitTestDB(t)
	})
}

func (s *repositoryTestSuite) TestPGRepository_FindAll() {
	ctx := context.Background()

	s.T().Run("should find all successfully", func(t *testing.T) {
		s.insertRepo(ctx, 5)

		names, err := s.agents.Repository().FindAll(ctx)

		assert.Equal(t, 5, len(names))
		assert.NoError(t, err)
	})

	s.T().Run("should return PostgresConnectionError if select fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)
		_, err := s.agents.Repository().FindAll(ctx)

		assert.True(t, errors.IsPostgresConnectionError(err))

		s.pg.InitTestDB(t)
	})
}

func (s *repositoryTestSuite) insertRepo(ctx context.Context, number int) {
	for i := 0; i < number; i++ {
		repo := &models.RepositoryModel{
			Name: "myRepository_" + strconv.Itoa(i),
		}

		_ = s.agents.Repository().Insert(ctx, repo)
	}
}
