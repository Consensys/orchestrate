// +build unit

package sessions

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestCreateSession_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewCreateSessionUseCase()

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJob()

		session := usecase.Execute(ctx, job)

		assert.Equal(t, job, session.Job)
		assert.NotEmpty(t, session.Cancel)
	})
}
