// +build unit
// +build !race
// +build !integration

package orm

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	storePG "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres/migrations"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

type jobTestSuite struct {
	suite.Suite
	orm ORM
	pg  *pgTestUtils.PGTestHelper
}

func TestPGJob(t *testing.T) {
	s := new(jobTestSuite)
	suite.Run(t, s)
}

func (s *jobTestSuite) SetupSuite() {
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *jobTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.orm = New()
}

func (s *jobTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *jobTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *jobTestSuite) TestORMJob_PGInsertOrUpdateJobSuccess() {
	job := testutils.FakeJob(0)
	pgdb := storePG.NewPGDB(s.pg.DB)

	s.T().Run("should insert Job model and its linked entities successfully", func(t *testing.T) {
		err := s.orm.InsertOrUpdateJob(context.Background(), pgdb, job)

		assert.Nil(t, err)
		assert.NotEmpty(t, job.ID)
		assert.NotEmpty(t, job.Transaction.ID)
		assert.NotEmpty(t, job.Schedule.ID)
		assert.NotEmpty(t, job.Logs[0].ID)
		assert.NotEmpty(t, job.UUID)
	})

	s.T().Run("should update Job model and its linked entities successfully", func(t *testing.T) {
		expectedJobData := testutils.FakeJob(0)
		job.TransactionID = nil
		job.Transaction.UUID = expectedJobData.Transaction.UUID
		job.Logs = append(job.Logs, &models.Log{Status: types.JobStatusCreated, Message: "created message"})
		err := s.orm.InsertOrUpdateJob(context.Background(), pgdb, job)

		assert.Nil(t, err)
		assert.Equal(t, expectedJobData.Transaction.UUID, job.Transaction.UUID)
		assert.Equal(t, 2, len(job.Logs))
		assert.NotEmpty(t, job.UUID)
	})
}

func (s *jobTestSuite) TestORMJob_MockInsertOrUpdateJobSuccess() {
	ctx := context.Background()
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	mockDBTX := mocks.NewMockTx(ctrl)

	mockScheduleDA := mocks.NewMockScheduleAgent(ctrl)
	mockLogDA := mocks.NewMockLogAgent(ctrl)
	mockTxDA := mocks.NewMockTransactionAgent(ctrl)
	mockJobDA := mocks.NewMockJobAgent(ctrl)

	mockDB.EXPECT().Begin().Return(mockDBTX, nil).AnyTimes()
	mockDBTX.EXPECT().Begin().Return(mockDBTX, nil).AnyTimes()
	mockDBTX.EXPECT().Transaction().Return(mockTxDA).AnyTimes()
	mockDBTX.EXPECT().Schedule().Return(mockScheduleDA).AnyTimes()
	mockDBTX.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDBTX.EXPECT().Log().Return(mockLogDA).AnyTimes()
	mockDBTX.EXPECT().Commit().Return(nil).AnyTimes()
	mockDBTX.EXPECT().Close().Return(nil).AnyTimes()

	s.T().Run("should insert Job model and its linked entities successfully", func(t *testing.T) {
		job := testutils.FakeJob(0)

		mockTxDA.EXPECT().Insert(gomock.Any(), gomock.Eq(job.Transaction)).Return(nil)
		mockScheduleDA.EXPECT().Insert(ctx, job.Schedule).Return(nil)
		mockJobDA.EXPECT().Insert(ctx, job).Return(nil)
		mockLogDA.EXPECT().Insert(ctx, job.Logs[0]).Return(nil)

		err := s.orm.InsertOrUpdateJob(ctx, mockDB, job)
		assert.Nil(t, err)
	})

	s.T().Run("should update Job model and its linked entities successfully", func(t *testing.T) {
		job := testutils.FakeJob(0)
		job.ID = 1
		job.Transaction.ID = 1
		job.Schedule.ID = 1 
		job.Logs[0].ID = 1 
		job.Logs = append(job.Logs, &models.Log{})
	
		mockTxDA.EXPECT().Update(gomock.Any(), gomock.Eq(job.Transaction)).Return(nil)
		mockJobDA.EXPECT().Update(ctx, job).Return(nil)
		mockLogDA.EXPECT().Insert(ctx, job.Logs[1]).Return(nil)
	
		err := s.orm.InsertOrUpdateJob(ctx, mockDB, job)
		assert.Nil(t, err)
	})
}
