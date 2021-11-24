package usecases

import (
	"context"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/store"
)

//go:generate mockgen -source=schedules.go -destination=mocks/schedules.go -package=mocks

/**
Schedule Use Cases
*/
type ScheduleUseCases interface {
	CreateSchedule() CreateScheduleUseCase
	GetSchedule() GetScheduleUseCase
	SearchSchedules() SearchSchedulesUseCase
}

type CreateScheduleUseCase interface {
	Execute(ctx context.Context, schedule *entities.Schedule, userInfo *multitenancy.UserInfo) (*entities.Schedule, error)
	WithDBTransaction(dbtx store.Tx) CreateScheduleUseCase
}

type GetScheduleUseCase interface {
	Execute(ctx context.Context, scheduleUUID string, userInfo *multitenancy.UserInfo) (*entities.Schedule, error)
}

type SearchSchedulesUseCase interface {
	Execute(ctx context.Context, userInfo *multitenancy.UserInfo) ([]*entities.Schedule, error)
}
