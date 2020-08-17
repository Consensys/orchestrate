// +build unit

package listeners

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/listeners/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs/mocks"
)

func TestJobsListener(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	searchJobUC := mocks.NewMockSearchJobsUseCase(ctrl)
	mockSessionManager := mocks2.NewMockSessionManager(ctrl)
	refreshInterval := 100 * time.Millisecond

	listener := NewJobsListener(refreshInterval, mockSessionManager, searchJobUC)

	t.Run("should add a new session for every retrieved job at every refresh", func(t *testing.T) {
		// We expect the listener to tick 5 times (initialCall + 4 times). We add 20ms to make sure we have a bit more time
		timeout := refreshInterval*5 + 20*time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), timeout)

		job0 := testutils.FakeJob()
		job1 := testutils.FakeJob()

		searchJobUC.EXPECT().Execute(ctx, gomock.Any(), []string{multitenancy.Wildcard}).Return([]*entities.Job{job0, job1}, nil).Times(5)
		mockSessionManager.EXPECT().AddSession(ctx, job0).Return(nil).Times(5)
		mockSessionManager.EXPECT().AddSession(ctx, job1).Return(nil).Times(5)

		cerr := listener.Listen(ctx)
		select {
		case err := <-cerr:
			assert.NoError(t, err)
		case <-ctx.Done():
		}
	})

	t.Run("should stop on context cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		cerr := listener.Listen(ctx)
		go sleepAnCancel(20*time.Millisecond, cancel)

		select {
		case err := <-cerr:
			assert.NoError(t, err)
		case <-ctx.Done():
			assert.Equal(t, "context canceled", ctx.Err().Error())
		}
	})

	t.Run("should fail if search jobs fail", func(t *testing.T) {
		ctx := context.Background() // context background making sure the context cannot be cancelled

		searchJobUC.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

		cerr := listener.Listen(ctx)
		select {
		case err := <-cerr:
			assert.Equal(t, errors.InternalError("failed to get latest pending jobs").ExtendComponent(jobsListenerComponent), err)
		}
	})

	t.Run("should fail if Add session fail", func(t *testing.T) {
		ctx := context.Background() // context background making sure the context cannot be cancelled
		job := testutils.FakeJob()

		searchJobUC.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return([]*entities.Job{job}, nil)
		mockSessionManager.EXPECT().AddSession(ctx, job).Return(fmt.Errorf("error"))

		cerr := listener.Listen(ctx)
		select {
		case err := <-cerr:
			assert.Equal(t, errors.InternalError("failed to add job session").ExtendComponent(jobsListenerComponent), err)
		}
	})
}

func sleepAnCancel(d time.Duration, cancel context.CancelFunc) {
	time.Sleep(d)
	cancel()
}
