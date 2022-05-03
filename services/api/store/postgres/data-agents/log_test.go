// +build !unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/errors"
	pgTestUtils "github.com/consensys/orchestrate/pkg/toolkit/database/postgres/testutils"
	"github.com/consensys/orchestrate/services/api/store/models/testutils"
	"github.com/consensys/orchestrate/services/api/store/postgres/migrations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
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

	job := testutils.FakeJobModel(0)
	err := s.agents.Schedule().Insert(ctx, job.Schedule)
	assert.NoError(s.T(), err)
	err = s.agents.Transaction().Insert(ctx, job.Transaction)
	assert.NoError(s.T(), err)

	chain := testutils.FakeChainModel()
	chain.UUID = job.ChainUUID
	err = s.agents.Chain().Insert(ctx, chain)
	assert.NoError(s.T(), err)

	err = s.agents.Job().Insert(ctx, job)
	assert.NoError(s.T(), err)

	s.T().Run("should insert model successfully", func(t *testing.T) {
		jobLog := testutils.FakeLog()
		jobLog.JobID = &job.ID
		err = s.agents.Log().Insert(ctx, jobLog)

		assert.NoError(t, err)
		assert.NotEmpty(t, jobLog.ID) // 2 because one is inserted when creating the job
	})

	s.T().Run("should insert model without UUID successfully", func(t *testing.T) {
		jobLog := testutils.FakeLog()
		jobLog.UUID = ""
		jobLog.JobID = &job.ID
		err = s.agents.Log().Insert(ctx, jobLog)

		assert.NoError(t, err)
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
		assert.True(t, errors.IsInternalError(err))
	})

	// We bring it back up
	s.pg.InitTestDB(s.T())
}
