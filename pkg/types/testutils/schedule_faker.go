package testutils

import (
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/gofrs/uuid"
)

func FakeSchedule() *entities.Schedule {
	scheduleUUID := uuid.Must(uuid.NewV4()).String()
	job := FakeJob()
	job.ScheduleUUID = scheduleUUID

	return &entities.Schedule{
		UUID: scheduleUUID,
		Jobs: []*entities.Job{job},
	}
}
