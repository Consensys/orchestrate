// +build unit

package listeners

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sentry/tx-sentry/use-cases/mocks"
)

func TestSessionManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	retryInterval := 1 * time.Second

	client := mock.NewMockOrchestrateClient(ctrl)
	retrySessionJobUC := mocks.NewMockRetrySessionJobUseCase(ctrl)

	sessionManager := NewSessionManager(client, retrySessionJobUC)

	t.Run("should retry session job successfully at every retry interval with latest child job", func(t *testing.T) {
		timeout := retryInterval*4 + 500*time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.RetryInterval = retryInterval
		childJobUUIDOne := "childJobUUIDOne"
		childJobUUIDTwo := "childJobUUIDTwo"

		parentJobResponse := testutils.FakeJobResponse()
		parentJobResponse.UUID = job.UUID
		client.EXPECT().SearchJob(gomock.Any(), &entities.JobFilters{
			ChainUUID:     job.ChainUUID,
			ParentJobUUID: job.UUID,
		}).Return([]*types.JobResponse{parentJobResponse}, nil)
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, parentJobResponse.UUID, 0).Return(childJobUUIDOne, nil)
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, childJobUUIDOne, 1).Return(childJobUUIDTwo, nil)
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, childJobUUIDTwo, 2).Return(childJobUUIDTwo, nil)
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, childJobUUIDTwo, 2).Return("", nil)
		client.EXPECT().UpdateJob(gomock.Any(), job.UUID, gomock.Any()).Return(nil, nil)

		sessionManager.Start(ctx, job)

		<-ctx.Done()
	})

	t.Run("should retry session job successfully at every retry interval with latest child job", func(t *testing.T) {
		timeout := retryInterval*4 + 500*time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.RetryInterval = retryInterval
		childJobUUIDOne := "childJobUUIDOne"
		childJobUUIDTwo := "childJobUUIDTwo"

		parentJobResponse := testutils.FakeJobResponse()
		parentJobResponse.UUID = job.UUID
		client.EXPECT().SearchJob(gomock.Any(), &entities.JobFilters{
			ChainUUID:     job.ChainUUID,
			ParentJobUUID: job.UUID,
		}).Return([]*types.JobResponse{parentJobResponse}, nil)
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, parentJobResponse.UUID, 0).Return(childJobUUIDOne, nil)
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, childJobUUIDOne, 1).Return(childJobUUIDTwo, nil)
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, childJobUUIDTwo, 2).Return(childJobUUIDTwo, nil)
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, childJobUUIDTwo, 2).Return("", nil)
		client.EXPECT().UpdateJob(gomock.Any(), job.UUID, gomock.Any()).Return(nil, nil)

		sessionManager.Start(ctx, job)

		<-ctx.Done()
	})

	t.Run("should not retry session job if max retries was exceeded", func(t *testing.T) {
		timeout := retryInterval*2 + 500*time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.RetryInterval = retryInterval

		parentJobResponse := testutils.FakeJobResponse()
		parentJobResponse.UUID = job.UUID
		for idx := 0; idx < 10; idx++ {
			parentJobResponse.Logs = append(parentJobResponse.Logs, &entities.Log{
				Status: utils.StatusResending,
			})
		}

		client.EXPECT().SearchJob(gomock.Any(), &entities.JobFilters{
			ChainUUID:     job.ChainUUID,
			ParentJobUUID: job.UUID,
		}).Return([]*types.JobResponse{parentJobResponse}, nil)

		sessionManager.Start(ctx, job)

		<-ctx.Done()
	})

	t.Run("should not retry session job in case of empty childJobUUID", func(t *testing.T) {
		timeout := retryInterval*2 + 500*time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.RetryInterval = retryInterval

		parentJobResponse := testutils.FakeJobResponse()
		parentJobResponse.UUID = job.UUID
		client.EXPECT().SearchJob(gomock.Any(), &entities.JobFilters{
			ChainUUID:     job.ChainUUID,
			ParentJobUUID: job.UUID,
		}).Return([]*types.JobResponse{parentJobResponse}, nil)
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, parentJobResponse.UUID, 0).Return("", nil)

		client.EXPECT().UpdateJob(gomock.Any(), job.UUID, gomock.Any()).Return(nil, nil)
		sessionManager.Start(ctx, job)

		<-ctx.Done()
	})

	t.Run("should retry session job successfully up to MaxRetry times", func(t *testing.T) {
		timeout := retryInterval*types.SentryMaxRetries + retryInterval*2 + 500*time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.RetryInterval = retryInterval
	
		parentJobResponse := testutils.FakeJobResponse()
		parentJobResponse.UUID = job.UUID
		client.EXPECT().SearchJob(gomock.Any(), &entities.JobFilters{
			ChainUUID:     job.ChainUUID,
			ParentJobUUID: job.UUID,
		}).Return([]*types.JobResponse{parentJobResponse}, nil)
	
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, parentJobResponse.UUID, gomock.Any()).
			Return(parentJobResponse.UUID, nil).Times(types.SentryMaxRetries)
	
		client.EXPECT().UpdateJob(gomock.Any(), job.UUID, gomock.Any()).Return(nil, nil)
		sessionManager.Start(ctx, job)
	
		<-ctx.Done()
	})

	t.Run("should retry with backoff if createChildJob fails", func(t *testing.T) {
		timeout := 2*retryInterval + 500*time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.RetryInterval = retryInterval
	
		parentJobResponse := testutils.FakeJobResponse()
		parentJobResponse.UUID = job.UUID
		client.EXPECT().SearchJob(gomock.Any(), &entities.JobFilters{
			ChainUUID:     job.ChainUUID,
			ParentJobUUID: job.UUID,
		}).Return([]*types.JobResponse{parentJobResponse}, nil)
	
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, parentJobResponse.UUID, 0).Return("", fmt.Errorf("error"))
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, parentJobResponse.UUID, gomock.Any()).Return(parentJobResponse.UUID, nil).AnyTimes()
		client.EXPECT().UpdateJob(gomock.Any(), job.UUID, gomock.Any()).Return(nil, nil)
		sessionManager.Start(ctx, job)
	
		<-ctx.Done()
	})

	t.Run("should do nothing if there is an active session for same job", func(t *testing.T) {
		timeout := retryInterval*2 + 600*time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.RetryInterval = retryInterval
	
		parentJobResponse := testutils.FakeJobResponse()
		parentJobResponse.UUID = job.UUID
		client.EXPECT().SearchJob(gomock.Any(), &entities.JobFilters{
			ChainUUID:     job.ChainUUID,
			ParentJobUUID: job.UUID,
		}).Return([]*types.JobResponse{parentJobResponse}, nil)
	
		// First session is added and startSessionUC is called
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, parentJobResponse.UUID, gomock.Any()).
			Return(parentJobResponse.UUID, nil)
		retrySessionJobUC.EXPECT().Execute(gomock.Any(), parentJobResponse.UUID, parentJobResponse.UUID, gomock.Any()).
			Return("", nil)
		sessionManager.Start(ctx, job)
	
		// Second call does not call startSessionUC
		ctx2, _ := context.WithTimeout(context.Background(), timeout)
		sessionManager.Start(ctx, job)
		client.EXPECT().UpdateJob(gomock.Any(), job.UUID, gomock.Any()).Return(nil, nil)
	
		<-ctx.Done()
		<-ctx2.Done()
	})

	t.Run("should do nothing if not retry interval is set", func(t *testing.T) {
		timeout := retryInterval + 500*time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.RetryInterval = 0

		// First session is added and startSessionUC is called
		sessionManager.Start(ctx, job)

		<-ctx.Done()
	})

	t.Run("should do nothing if job was already resent", func(t *testing.T) {
		timeout := retryInterval + 500*time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.RetryInterval = retryInterval
		job.InternalData.HasBeenRetried = true

		sessionManager.Start(ctx, job)

		<-ctx.Done()
	})
}
