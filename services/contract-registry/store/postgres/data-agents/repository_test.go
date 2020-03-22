// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/postgres/migrations"
)

type repositoryTestSuite struct {
	suite.Suite
	dataagent *PGRepository
	pg        *pgTestUtils.PGTestHelper
}

func TestPGRepository(t *testing.T) {
	s := new(repositoryTestSuite)
	suite.Run(t, s)
}

func (s *repositoryTestSuite) SetupSuite() {
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *repositoryTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())
	s.dataagent = NewPGRepository(s.pg.DB)
}

func (s *repositoryTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *repositoryTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *repositoryTestSuite) TestPGRepository_SelectOrInsertTx() {
	s.T().Run("should insert model successfully", func(t *testing.T) {
		repo := &models.RepositoryModel{
			Name: "myRepository",
		}
		err := s.dataagent.SelectOrInsert(context.Background(), repo)

		assert.Nil(t, err)
		assert.Equal(t, 1, repo.ID)
	})

	s.T().Run("should insert model successfully in TX", func(t *testing.T) {
		tx, _ := s.pg.DB.Begin()
		ctx := postgres.WithTx(context.Background(), tx)

		repo := &models.RepositoryModel{
			Name: "myRepository",
		}
		err := s.dataagent.SelectOrInsert(ctx, repo)
		_ = tx.Commit()

		assert.Nil(t, err)
		assert.NotEmpty(t, repo.ID)
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		repo := &models.RepositoryModel{
			Name: "myRepository",
		}
		err := s.dataagent.SelectOrInsert(context.Background(), repo)

		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *repositoryTestSuite) TestPGRepository_FindAll() {
	s.T().Run("should find all successfully", func(t *testing.T) {
		s.insertRepo(5)

		names, err := s.dataagent.FindAll(context.Background())

		assert.Equal(t, 5, len(names))
		assert.Nil(t, err)
	})

	s.T().Run("should return PostgresConnectionError if select fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)
		_, err := s.dataagent.FindAll(context.Background())

		assert.True(t, errors.IsPostgresConnectionError(err))

		s.pg.InitTestDB(t)
	})
}

func (s *repositoryTestSuite) insertRepo(number int) {
	for i := 0; i < number; i++ {
		repo := &models.RepositoryModel{
			Name: "myRepository_" + strconv.Itoa(i),
		}

		_ = s.dataagent.SelectOrInsert(context.Background(), repo)
	}
}
