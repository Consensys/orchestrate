package builder

import (
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/schedules"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

type scheduleUseCases struct {
	createSchedule  usecases.CreateScheduleUseCase
	getSchedule     usecases.GetScheduleUseCase
	searchSchedules usecases.SearchSchedulesUseCase
}

func newScheduleUseCases(db store.DB) *scheduleUseCases {
	return &scheduleUseCases{
		createSchedule:  schedules.NewCreateScheduleUseCase(db),
		getSchedule:     schedules.NewGetScheduleUseCase(db),
		searchSchedules: schedules.NewSearchSchedulesUseCase(db),
	}
}

func (u *scheduleUseCases) CreateSchedule() usecases.CreateScheduleUseCase {
	return u.createSchedule
}

func (u *scheduleUseCases) GetSchedule() usecases.GetScheduleUseCase {
	return u.getSchedule
}

func (u *scheduleUseCases) SearchSchedules() usecases.SearchSchedulesUseCase {
	return u.searchSchedules
}
