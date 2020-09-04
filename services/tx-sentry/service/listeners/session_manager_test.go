// +build unit

package listeners

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-sentry/tx-sentry/use-cases/mocks"
	"testing"
	"time"
)

func TestSessionManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	retryInterval := 100 * time.Millisecond
	childJobUUID := "childJobUUID"

	mockCreateChilUC := mocks.NewMockCreateChildJobUseCase(ctrl)

	sessionManager := NewSessionManager(mockCreateChilUC)

	t.Run("should create a new child job successfully at every retry interval", func(t *testing.T) {
		// We expect to tick 5 times. We add 20ms to make sure we have a bit more time
		timeout := retryInterval*5 + 20*time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.RetryInterval = retryInterval

		mockCreateChilUC.EXPECT().Execute(ctx, job).Return(childJobUUID, nil).Times(5)

		sessionManager.Start(ctx, job)

		<-ctx.Done()
	})

	t.Run("should do nothing if session already exists", func(t *testing.T) {
		timeout := retryInterval + 20*time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.RetryInterval = retryInterval

		// First session is added and startSessionUC is called
		mockCreateChilUC.EXPECT().Execute(ctx, job).Return(childJobUUID, nil)
		sessionManager.Start(ctx, job)

		// Second call does not call startSessionUC
		sessionManager.Start(ctx, job)

		<-ctx.Done()
	})

	t.Run("should do nothing if session is not a parent job", func(t *testing.T) {
		timeout := 20 * time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.ParentJobUUID = "I am a child job"

		sessionManager.Start(ctx, job)
		<-ctx.Done()
	})

	t.Run("should stop the session if no child is created but no error", func(t *testing.T) {
		timeout := 1 * time.Second
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.RetryInterval = retryInterval

		mockCreateChilUC.EXPECT().Execute(ctx, job).Return("", nil)

		sessionManager.Start(ctx, job)

		<-ctx.Done()
	})

	t.Run("should retry with backoff if createChildJob fails", func(t *testing.T) {
		timeout := 1 * time.Second
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		job := testutils.FakeJob()
		job.InternalData.RetryInterval = retryInterval

		mockCreateChilUC.EXPECT().Execute(ctx, job).Return("", fmt.Errorf("error"))
		mockCreateChilUC.EXPECT().Execute(ctx, job).Return("", nil)

		sessionManager.Start(ctx, job)

		<-ctx.Done()
	})
}
