package schedules

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	tsorm "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/orm"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
)

//go:generate mockgen -source=get_schedule.go -destination=mocks/get_schedule.go -package=mocks

const getScheduleComponent = "use-cases.get-schedule"

type GetScheduleUseCase interface {
	Execute(ctx context.Context, scheduleUUID, tenantID string) (*types.ScheduleResponse, error)
}

// getScheduleUseCase is a use case to get a schedule
type getScheduleUseCase struct {
	db  store.DB
	orm tsorm.ORM
}

// NewGetScheduleUseCase creates a new GetScheduleUseCase
func NewGetScheduleUseCase(db store.DB, orm tsorm.ORM) GetScheduleUseCase {
	return &getScheduleUseCase{
		db:  db,
		orm: orm,
	}
}

// Execute gets a schedule
func (uc *getScheduleUseCase) Execute(ctx context.Context, scheduleUUID, tenantID string) (*types.ScheduleResponse, error) {
	log.WithContext(ctx).
		WithField("schedule_uuid", scheduleUUID).
		Debug("getting schedule")

	schedule, err := uc.orm.FetchScheduleByUUID(ctx, uc.db, scheduleUUID, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getScheduleComponent)
	}

	log.WithContext(ctx).
		WithField("schedule_uuid", schedule.UUID).
		Info("schedule found successfully")
	return utils.FormatScheduleResponse(schedule), nil
}
