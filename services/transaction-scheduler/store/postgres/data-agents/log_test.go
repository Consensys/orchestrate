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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/migrations"
)

type logTestSuite struct {
	suite.Suite
	agents *PGAgents
	pg     *pgTestUtils.PGTestHelper
}

func TestPGLog(t *testing.T) {
	s := new(logTestSuite)
	suite.Run(t, s)
}

func (s *logTestSuite) SetupSuite() {
	s.pg , _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *logTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.agents = New(s.pg.DB)
}

func (s *logTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *logTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *logTestSuite) TestPGLog_Insert() {
	ctx := context.Background()

	job := testutils.FakeJob(0)
	err := s.agents.Schedule().Insert(ctx, job.Schedule)
	assert.Nil(s.T(), err)
	err = s.agents.Transaction().Insert(ctx, job.Transaction)
	assert.Nil(s.T(), err)
	err = s.agents.Job().Insert(ctx, job)
	assert.Nil(s.T(), err)

	s.T().Run("should insert model successfully", func(t *testing.T) {
		jobLog := testutils.FakeLog()
		jobLog.JobID = &job.ID
		err = s.agents.Log().Insert(ctx, jobLog)

		assert.Nil(t, err)
		assert.NotEmpty(t, jobLog.ID) // 2 because one is inserted when creating the job
	})

	s.T().Run("should insert model without UUID successfully", func(t *testing.T) {
		jobLog := testutils.FakeLog()
		jobLog.UUID = ""
		jobLog.JobID = &job.ID
		err = s.agents.Log().Insert(ctx, jobLog)

		assert.Nil(t, err)
		assert.NotEmpty(t, jobLog.ID) // 2 because one is inserted when creating the job
	})
}

func (s *logTestSuite) TestPGLog_ConnectionErr() {
	ctx := context.Background()

	// We drop the DB to make the test fail
	s.pg.DropTestDB(s.T())
	job := testutils.FakeLog()
	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		err := s.agents.Log().Insert(ctx, job)
		assert.True(t, errors.IsPostgresConnectionError(err))
	})

	// We bring it back up
	s.pg.InitTestDB(s.T())
}
