// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/migrations"
	"testing"
)

type logTestSuite struct {
	suite.Suite
	dataagent  *PGLog
	jobDA      *PGJob
	scheduleDA *PGSchedule
	pg         *pgTestUtils.PGTestHelper
}

func TestPGLog(t *testing.T) {
	s := new(logTestSuite)
	suite.Run(t, s)
}

func (s *logTestSuite) SetupSuite() {
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *logTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())
	s.jobDA = NewPGJob(s.pg.DB)
	s.scheduleDA = NewPGSchedule(s.pg.DB)
	s.dataagent = NewPGLog(s.pg.DB)
}

func (s *logTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *logTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *logTestSuite) TestPGLog_Insert() {
	ctx := context.Background()
	schedule := testutils.FakeSchedule()
	_ = s.scheduleDA.Insert(ctx, schedule)

	job := testutils.FakeJob(schedule.ID)
	_ = s.jobDA.Insert(ctx, job)

	s.T().Run("should insert model successfully", func(t *testing.T) {
		log := testutils.FakeLog(job.ID)
		err := s.dataagent.Insert(context.Background(), log)

		assert.Nil(t, err)
		assert.NotNil(t, log.UUID)
		assert.Equal(t, 2, log.ID) // 2 because one is inserted when creating the job
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		log := testutils.FakeLog(job.ID)
		err := s.dataagent.Insert(context.Background(), log)
		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}
