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

type scheduleTestSuite struct {
	suite.Suite
	dataagent *PGSchedule
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

		err := s.insertSchedule()
		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *scheduleTestSuite) insertSchedule() error {
	return s.dataagent.Insert(context.Background(), testutils.FakeSchedule())
}
