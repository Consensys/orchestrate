package schedules

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
)

//go:generate mockgen -source=get_schedules.go -destination=mocks/get_schedules.go -package=mocks

const getSchedulesComponent = "use-cases.get-schedules"

type GetSchedulesUseCase interface {
	Execute(ctx context.Context, tenantID string) ([]*entities.Schedule, error)
}

// getScheduleUseCase is a use case to get a schedule
type getSchedulesUseCase struct {
	db store.DB
}

// NewGetScheduleUseCase creates a new GetScheduleUseCase
func NewGetSchedulesUseCase(db store.DB) GetSchedulesUseCase {
	return &getSchedulesUseCase{
		db: db,
	}
}

// Execute gets a schedule
func (uc *getSchedulesUseCase) Execute(ctx context.Context, tenantID string) ([]*entities.Schedule, error) {
	log.WithContext(ctx).
		WithField("schedule.tenantID", tenantID).
		Debug("getting schedule")

	scheduleModels, err := fetchAllSchedule(ctx, uc.db, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getSchedulesComponent)
	}

	log.WithContext(ctx).
		WithField("schedule.tenantID", tenantID).
		Info("schedule found successfully")

	resp := []*entities.Schedule{}
	for _, s := range scheduleModels {
		resp = append(resp, parsers.NewScheduleEntityFromModels(s))
	}

	return resp, nil
}

func fetchAllSchedule(ctx context.Context, db store.DB, tenantID string) ([]*models.Schedule, error) {
	schedules, err := db.Schedule().FindAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	for idx, schedule := range schedules {
		for jdx, job := range schedule.Jobs {
			schedules[idx].Jobs[jdx], err = db.Job().FindOneByUUID(ctx, job.UUID, tenantID)
			if err != nil {
				return nil, err
			}
		}
	}

	return schedules, nil
}
