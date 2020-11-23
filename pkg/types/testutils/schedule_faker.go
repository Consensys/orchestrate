package testutils

import (
	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
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
