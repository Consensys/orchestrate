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

type jobTestSuite struct {
	suite.Suite
	dataagent  *PGJob
	scheduleDA *PGSchedule
	pg         *pgTestUtils.PGTestHelper
}

func TestPGJob(t *testing.T) {
	s := new(jobTestSuite)
	suite.Run(t, s)
}

func (s *jobTestSuite) SetupSuite() {
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *jobTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())
	s.scheduleDA = NewPGSchedule(s.pg.DB)
	s.dataagent = NewPGJob(s.pg.DB)
}

func (s *jobTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *jobTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *jobTestSuite) TestPGJob_Insert() {
	ctx := context.Background()
	schedule := testutils.FakeSchedule()
	_ = s.scheduleDA.Insert(ctx, schedule)

	s.T().Run("should insert model successfully", func(t *testing.T) {
		job := testutils.FakeJob(schedule.ID)
		err := s.dataagent.Insert(context.Background(), job)

		assert.Nil(t, err)
		assert.Equal(t, 1, job.ID)
		assert.NotNil(t, job.UUID)
		assert.Equal(t, 1, job.TransactionID)
		assert.NotNil(t, job.Transaction.UUID)
		assert.Equal(t, job.Transaction.ID, job.TransactionID)
		assert.Equal(t, 1, job.Logs[0].ID)
		assert.NotNil(t, job.Logs[0].UUID)
		assert.Equal(t, job.Logs[0].JobID, job.ID)
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		job := testutils.FakeJob(schedule.ID)
		err := s.dataagent.Insert(context.Background(), job)
		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}
