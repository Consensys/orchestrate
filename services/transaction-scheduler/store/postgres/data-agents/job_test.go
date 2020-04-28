// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
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
	logDA      *PGLog
	pg         *pgTestUtils.PGTestHelper
}

func TestPGJob(t *testing.T) {
	s := new(jobTestSuite)
	suite.Run(t, s)
}

func (s *jobTestSuite) SetupSuite() {
	s.pg , _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *jobTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.scheduleDA = NewPGSchedule(s.pg.DB)
	s.logDA = NewPGLog(s.pg.DB)
	s.dataagent = NewPGJob(s.pg.DB)
}

func (s *jobTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
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

func (s *jobTestSuite) TestPGJob_FindOneByUUID() {
	ctx := context.Background()
	tenantID := "tenantID"
	schedule := testutils.FakeSchedule()
	_ = s.scheduleDA.Insert(ctx, schedule)
	job := testutils.FakeJob(schedule.ID)
	_ = s.dataagent.Insert(context.Background(), job)
	job.Logs[0].JobID = job.ID
	_ = s.logDA.Insert(context.Background(), job.Logs[0])

	s.T().Run("should get model successfully as tenant", func(t *testing.T) {
		jobRetrieved, err := s.dataagent.FindOneByUUID(ctx, job.UUID, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, job.ID, jobRetrieved.ID)
		assert.Equal(t, job.UUID, jobRetrieved.UUID)
		assert.Equal(t, job.Transaction, jobRetrieved.Transaction)
		assert.Equal(t, job.Logs, jobRetrieved.Logs)
		assert.Equal(t, schedule.UUID, jobRetrieved.Schedule.UUID)
	})

	s.T().Run("should get model successfully as admin", func(t *testing.T) {
		jobRetrieved, err := s.dataagent.FindOneByUUID(ctx, job.UUID, multitenancy.DefaultTenantIDName)

		assert.Nil(t, err)
		assert.Equal(t, job.ID, jobRetrieved.ID)
	})

	s.T().Run("should return NotFoundError if select fails", func(t *testing.T) {
		_, err := s.dataagent.FindOneByUUID(ctx, "b6fe7a2a-1a4d-49ca-99d8-8a34aa495ef0", tenantID)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		_, err := s.dataagent.FindOneByUUID(ctx, job.UUID, tenantID)
		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}
