// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/postgres/migrations"
)

type jobTestSuite struct {
	suite.Suite
	agents   *PGAgents
	pg       *pgTestUtils.PGTestHelper
	tenantID string
}

func TestPGJob(t *testing.T) {
	s := new(jobTestSuite)
	suite.Run(t, s)
}

func (s *jobTestSuite) SetupSuite() {
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.tenantID = "tenantID"
	s.pg.InitTestDB(s.T())
}

func (s *jobTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.agents = New(s.pg.DB)
}

func (s *jobTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *jobTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *jobTestSuite) TestPGJob_Insert() {
	ctx := context.Background()

	s.T().Run("should insert model successfully", func(t *testing.T) {
		job := testutils.FakeJobModel(0)
		err := insertJob(ctx, s.agents, job)
		assert.NoError(s.T(), err)

		assert.NoError(t, err)
		assert.NotEmpty(t, job.ID)
		assert.NotEmpty(t, job.Transaction.ID)
		assert.NotEmpty(t, job.Schedule.ID)
	})

	s.T().Run("should insert model without UUID successfully", func(t *testing.T) {
		job := testutils.FakeJobModel(0)
		job.UUID = ""
		err := insertJob(ctx, s.agents, job)
		assert.NoError(s.T(), err)

		assert.NoError(t, err)
		assert.NotEmpty(t, job.ID)
		assert.NotEmpty(t, job.Transaction.ID)
		assert.NotEmpty(t, job.Schedule.ID)
	})

	s.T().Run("should update model successfully", func(t *testing.T) {
		job := testutils.FakeJobModel(0)
		err := insertJob(ctx, s.agents, job)
		assert.NoError(s.T(), err)

		assert.NoError(t, err)
		assert.NotEmpty(t, job.ID)
		assert.NotEmpty(t, job.Transaction.ID)
		assert.NotEmpty(t, job.Schedule.ID)
	})
}

func (s *jobTestSuite) TestPGJob_Update() {
	ctx := context.Background()
	job := testutils.FakeJobModel(0)
	err := insertJob(ctx, s.agents, job)
	assert.NoError(s.T(), err)

	s.T().Run("should update model successfully", func(t *testing.T) {
		newTx := testutils.FakeTransaction()
		newSchedule := testutils.FakeSchedule("_")
		err = s.agents.Transaction().Insert(ctx, newTx)
		assert.NoError(t, err)
		err = s.agents.Schedule().Insert(ctx, newSchedule)
		assert.NoError(t, err)

		job.ScheduleID = &newSchedule.ID
		job.TransactionID = &newTx.ID
		err = s.agents.Job().Update(ctx, job)
		assert.NoError(t, err)
		assert.Equal(t, *job.TransactionID, newTx.ID)
		assert.Equal(t, *job.ScheduleID, newSchedule.ID)
	})

	s.T().Run("should fail to update job with missing ID", func(t *testing.T) {
		job.ID = 0
		err = s.agents.Job().Update(ctx, job)
		assert.True(t, errors.IsInvalidArgError(err))
	})
}

func (s *jobTestSuite) TestPGJob_FindOneByUUID() {
	ctx := context.Background()
	job := testutils.FakeJobModel(0)
	job.NextJobUUID = uuid.Must(uuid.NewV4()).String()
	job.Logs = append(job.Logs, &models.Log{UUID: uuid.Must(uuid.NewV4()).String(), Status: entities.StatusStarted, Message: "created message"})
	job.Schedule.TenantID = s.tenantID

	err := insertJob(ctx, s.agents, job)
	assert.NoError(s.T(), err)

	s.T().Run("should get model successfully with sorted logs", func(t *testing.T) {
		jobRetrieved, err := s.agents.Job().FindOneByUUID(ctx, job.UUID, []string{multitenancy.Wildcard})

		assert.NoError(t, err)
		assert.NotEmpty(t, jobRetrieved.ID)
		assert.Equal(t, job.UUID, jobRetrieved.UUID)
		assert.Equal(t, job.NextJobUUID, jobRetrieved.NextJobUUID)
		assert.Equal(t, job.Transaction.UUID, jobRetrieved.Transaction.UUID)
		assert.NotEmpty(t, jobRetrieved.Transaction.ID)
		assert.Equal(t, job.Logs[0].UUID, jobRetrieved.Logs[0].UUID)
		assert.Equal(t, job.Logs[1].UUID, jobRetrieved.Logs[1].UUID)
		assert.NotEmpty(t, jobRetrieved.Logs[0].ID)
		assert.Equal(t, job.Schedule.UUID, jobRetrieved.Schedule.UUID)
		assert.Equal(t, job.Schedule.TenantID, jobRetrieved.Schedule.TenantID)
		assert.NotEmpty(t, jobRetrieved.Schedule.ID)
	})

	s.T().Run("should get model successfully as tenant", func(t *testing.T) {
		jobRetrieved, err := s.agents.Job().FindOneByUUID(ctx, job.UUID, []string{s.tenantID})

		assert.NoError(t, err)
		assert.NotEmpty(t, jobRetrieved.ID)
	})

	s.T().Run("should return NotFoundError if select fails", func(t *testing.T) {
		_, err := s.agents.Job().FindOneByUUID(ctx, "b6fe7a2a-1a4d-49ca-99d8-8a34aa495ef0", []string{s.tenantID})
		assert.True(t, errors.IsNotFoundError(err))
	})
}

func (s *jobTestSuite) TestPGJob_LockOneByUUID() {
	ctx := context.Background()
	job := testutils.FakeJobModel(0)
	job.Logs = append(job.Logs, &models.Log{UUID: uuid.Must(uuid.NewV4()).String(), Status: entities.StatusStarted, Message: "created message"})
	job.Schedule.TenantID = s.tenantID
	err := insertJob(ctx, s.agents, job)
	assert.NoError(s.T(), err)

	s.T().Run("should lock successfully", func(t *testing.T) {
		dbtx0, err := s.pg.DB.Begin()
		assert.NoError(t, err)
		dbtx1, err := s.pg.DB.Begin()
		assert.NoError(t, err)
		newPGJob0 := NewPGJob(dbtx0)
		newPGJob1 := NewPGJob(dbtx1)

		waitChannel := make(chan string)
		err = newPGJob1.LockOneByUUID(ctx, job.UUID)
		assert.NoError(t, err)
		go func() {
			time.Sleep(2000 * time.Millisecond)
			waitChannel <- "job1"

			err = dbtx1.Commit()
			assert.NoError(t, err)
		}()

		go func() {
			err = newPGJob0.LockOneByUUID(ctx, job.UUID)
			assert.NoError(t, err)

			err = dbtx0.Commit()
			assert.NoError(t, err)

			waitChannel <- "job0"
		}()

		firstJob := <-waitChannel
		assert.Equal(t, firstJob, "job1")

		secondJob := <-waitChannel
		assert.Equal(t, secondJob, "job0")
	})
}

func (s *jobTestSuite) TestPGJob_Search() {
	ctx := context.Background()

	// job0 is the parent of a random job "parentJobUUID"
	job0 := testutils.FakeJobModel(0)
	job0.Logs = append(job0.Logs, &models.Log{UUID: uuid.Must(uuid.NewV4()).String(), Status: entities.StatusStarted, Message: "created message"})
	txHashOne := common.HexToHash("0x1")
	job0.Transaction.Hash = txHashOne.String()
	job0.Schedule.TenantID = s.tenantID
	job0.IsParent = true
	err := insertJob(ctx, s.agents, job0)
	assert.NoError(s.T(), err)

	// Job1 is the child of job0
	job1 := testutils.FakeJobModel(0)
	txHashTwo := common.HexToHash("0x2")
	job1.ChainUUID = job0.ChainUUID
	job1.Transaction.Hash = txHashTwo.String()
	job1.Schedule.TenantID = s.tenantID
	job1.InternalData.ParentJobUUID = job0.UUID
	job1.Logs[0].Status = entities.StatusPending
	err = insertJob(ctx, s.agents, job1)
	assert.NoError(s.T(), err)

	s.T().Run("should find model successfully", func(t *testing.T) {
		filters := &entities.JobFilters{
			TxHashes:  []string{txHashOne.String()},
			ChainUUID: job0.ChainUUID,
		}

		retrievedJobs, err := s.agents.Job().Search(ctx, filters, []string{s.tenantID})

		assert.NoError(t, err)
		assert.NotEmpty(t, retrievedJobs[0].ID)
		assert.Equal(t, job0.UUID, retrievedJobs[0].UUID)
		assert.Equal(t, job0.Transaction.UUID, retrievedJobs[0].Transaction.UUID)
		assert.Equal(t, txHashOne.String(), retrievedJobs[0].Transaction.Hash)
		assert.Equal(t, len(job0.Logs), len(retrievedJobs[0].Logs))
		// Verify order
		assert.Equal(t, job0.Logs[0].UUID, retrievedJobs[0].Logs[0].UUID)
		assert.Equal(t, job0.Logs[1].UUID, retrievedJobs[0].Logs[1].UUID)
	})

	s.T().Run("should not find any model by txHashes", func(t *testing.T) {
		filters := &entities.JobFilters{
			TxHashes:  []string{"0x3"},
			ChainUUID: job0.ChainUUID,
		}

		retrievedJobs, err := s.agents.Job().Search(ctx, filters, []string{s.tenantID})
		assert.NoError(t, err)
		assert.Empty(t, retrievedJobs)
	})

	s.T().Run("should not find any model by chainUUID", func(t *testing.T) {
		filters := &entities.JobFilters{
			TxHashes:  []string{txHashOne.String()},
			ChainUUID: uuid.Must(uuid.NewV4()).String(),
		}

		retrievedJobs, err := s.agents.Job().Search(ctx, filters, []string{s.tenantID})
		assert.NoError(t, err)
		assert.Empty(t, retrievedJobs)
	})

	s.T().Run("should find models successfully by parentJobUUID", func(t *testing.T) {
		// job0 is the parent so we retrieve the parent and all the children
		filters := &entities.JobFilters{
			ParentJobUUID: job0.UUID,
		}

		retrievedJobs, err := s.agents.Job().Search(ctx, filters, []string{s.tenantID})
		assert.NoError(t, err)
		assert.Len(t, retrievedJobs, 2)

		assert.Equal(t, retrievedJobs[0].UUID, job0.UUID)
		assert.Equal(t, retrievedJobs[1].InternalData.ParentJobUUID, job0.UUID)
		assert.Equal(t, retrievedJobs[1].UUID, job1.UUID)
	})

	s.T().Run("should find models successfully if OnlyParents is true", func(t *testing.T) {
		// job0 is the parent so we retrieve the parent and all the children
		filters := &entities.JobFilters{
			OnlyParents: true,
		}

		retrievedJobs, err := s.agents.Job().Search(ctx, filters, []string{s.tenantID})

		assert.NoError(t, err)
		assert.Len(t, retrievedJobs, 1)
		assert.Equal(t, retrievedJobs[0].UUID, job0.UUID)
	})

	s.T().Run("should find every inserted model successfully", func(t *testing.T) {
		filters := &entities.JobFilters{}
		retrievedJobs, err := s.agents.Job().Search(ctx, filters, []string{s.tenantID})

		assert.NoError(t, err)
		assert.Equal(t, len(retrievedJobs), 2)
	})
}

func (s *jobTestSuite) TestPGJob_ConnectionErr() {
	ctx := context.Background()

	// We drop the DB to make the test fail
	s.pg.DropTestDB(s.T())
	job := testutils.FakeJobModel(0)
	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		err := s.agents.Job().Insert(ctx, job)
		assert.True(t, errors.IsPostgresConnectionError(err))
	})

	s.T().Run("should return PostgresConnectionError if update fails", func(t *testing.T) {
		job.ID = 1
		err := s.agents.Job().Update(ctx, job)
		assert.True(t, errors.IsPostgresConnectionError(err))
	})

	s.T().Run("should return PostgresConnectionError if update fails", func(t *testing.T) {
		_, err := s.agents.Job().FindOneByUUID(ctx, job.UUID, []string{job.Schedule.TenantID})
		assert.True(t, errors.IsPostgresConnectionError(err))
	})

	// We bring it back up
	s.pg.InitTestDB(s.T())
}

/**
Persist Job entity and its related entities
*/
func insertJob(ctx context.Context, agents *PGAgents, job *models.Job) error {
	if _, err := agents.Chain().FindOneByUUID(ctx, job.ChainUUID, []string{}); errors.IsNotFoundError(err) {
		chain := testutils.FakeChainModel()
		chain.UUID = job.ChainUUID
		if err := agents.Chain().Insert(ctx, chain); err != nil {
			return err
		}
	}

	if job.Schedule != nil {
		if err := agents.Schedule().Insert(ctx, job.Schedule); err != nil {
			return err
		}
	}

	if job.Transaction != nil {
		if err := agents.Transaction().Insert(ctx, job.Transaction); err != nil {
			return err
		}
	}

	if err := agents.Job().Insert(ctx, job); err != nil {
		return err
	}

	for idx := range job.Logs {
		job.Logs[idx].JobID = &job.ID
		if err := agents.Log().Insert(ctx, job.Logs[idx]); err != nil {
			return err
		}
	}

	return nil
}
