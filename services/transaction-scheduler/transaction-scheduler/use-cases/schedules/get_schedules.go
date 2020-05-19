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

//go:generate mockgen -source=get_schedules.go -destination=mocks/get_schedules.go -package=mocks

const getSchedulesComponent = "use-cases.get-schedules"

type GetSchedulesUseCase interface {
	Execute(ctx context.Context, tenantID string) ([]*types.ScheduleResponse, error)
}

// getScheduleUseCase is a use case to get a schedule
type getSchedulesUseCase struct {
	db  store.DB
	orm tsorm.ORM
}

// NewGetScheduleUseCase creates a new GetScheduleUseCase
func NewGetSchedulesUseCase(db store.DB, orm tsorm.ORM) GetSchedulesUseCase {
	return &getSchedulesUseCase{
		db:  db,
		orm: orm,
	}
}

// Execute gets a schedule
func (uc *getSchedulesUseCase) Execute(ctx context.Context, tenantID string) ([]*types.ScheduleResponse, error) {
	log.WithContext(ctx).
		WithField("schedule.tenantID", tenantID).
		Debug("getting schedule")

	schedules, err := uc.orm.FetchAllSchedules(ctx, uc.db, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getSchedulesComponent)
	}

	log.WithContext(ctx).
		WithField("schedule.tenantID", tenantID).
		Info("schedule found successfully")

	resp := []*types.ScheduleResponse{}
	for _, s := range schedules {
		resp = append(resp, utils.FormatScheduleResponse(s))
	}
	return resp, nil
}
