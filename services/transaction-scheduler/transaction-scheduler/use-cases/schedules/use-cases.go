package schedules

type UseCases interface {
	CreateSchedule() CreateScheduleUseCase
	GetSchedule() GetScheduleUseCase
	GetSchedules() GetSchedulesUseCase
}
