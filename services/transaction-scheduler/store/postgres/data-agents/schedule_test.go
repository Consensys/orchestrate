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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/migrations"
	"testing"
)

type scheduleTestSuite struct {
	suite.Suite
	dataagent *PGSchedule
	jobDA     *PGJob
	pg        *pgTestUtils.PGTestHelper
}

func TestPGSchedule(t *testing.T) {
	s := new(scheduleTestSuite)
	suite.Run(t, s)
}

func (s *scheduleTestSuite) SetupSuite() {
	s.pg = pgTestUtils.NewPGTestHelper(migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *scheduleTestSuite) SetupTest() {
	s.pg.Upgrade(s.T())
	s.jobDA = NewPGJob(s.pg.DB)
	s.dataagent = NewPGSchedule(s.pg.DB)
}

func (s *scheduleTestSuite) TearDownTest() {
	s.pg.Downgrade(s.T())
}

func (s *scheduleTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *scheduleTestSuite) TestPGTransactionRequest_Insert() {
	s.T().Run("should insert model successfully", func(t *testing.T) {
		schedule := testutils.FakeSchedule()
		err := s.dataagent.Insert(context.Background(), schedule)

		assert.Nil(t, err)
		assert.NotNil(t, schedule.UUID)
		assert.Equal(t, schedule.ID, 1)
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		schedule := testutils.FakeSchedule()
		err := s.dataagent.Insert(context.Background(), schedule)
		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *scheduleTestSuite) TestPGJob_FindOneByUUID() {
	ctx := context.Background()
	tenantID := "tenantID"
	schedule := testutils.FakeSchedule()
	schedule.TenantID = tenantID
	_ = s.dataagent.Insert(ctx, schedule)

	schedule.Jobs[0].ScheduleID = schedule.ID
	_ = s.jobDA.Insert(ctx, schedule.Jobs[0])

	s.T().Run("should get model successfully as tenant", func(t *testing.T) {
		scheduleRetrieved, err := s.dataagent.FindOneByUUID(ctx, schedule.UUID, tenantID)

		assert.Nil(t, err)
		assert.Equal(t, schedule.ID, scheduleRetrieved.ID)
		assert.Equal(t, schedule.UUID, scheduleRetrieved.UUID)
		assert.Equal(t, schedule.ChainUUID, scheduleRetrieved.ChainUUID)
		assert.Equal(t, schedule.CreatedAt, scheduleRetrieved.CreatedAt)
		assert.Equal(t, schedule.Jobs[0].UUID, scheduleRetrieved.Jobs[0].UUID)
	})

	s.T().Run("should get model successfully as admin", func(t *testing.T) {
		jobRetrieved, err := s.dataagent.FindOneByUUID(ctx, schedule.UUID, multitenancy.DefaultTenantIDName)

		assert.Nil(t, err)
		assert.Equal(t, schedule.Jobs[0].ID, jobRetrieved.ID)
	})

	s.T().Run("should return NotFoundError if select fails", func(t *testing.T) {
		_, err := s.dataagent.FindOneByUUID(ctx, "b6fe7a2a-1a4d-49ca-99d8-8a34aa495ef0", tenantID)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		_, err := s.dataagent.FindOneByUUID(ctx, schedule.UUID, tenantID)
		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}
