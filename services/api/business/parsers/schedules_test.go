// +build unit

package parsers

import (
	testutils2 "github.com/ConsenSys/orchestrate/pkg/types/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/orchestrate/services/api/store/models/testutils"

	"github.com/ConsenSys/orchestrate/pkg/encoding/json"
)

func TestParsersSchedule_NewModelFromEntity(t *testing.T) {
	scheduleEntity := testutils2.FakeSchedule()
	scheduleModel := NewScheduleModelFromEntities(scheduleEntity)
	finalScheduleEntity := NewScheduleEntityFromModels(scheduleModel)

	expectedJSON, _ := json.Marshal(scheduleEntity)
	actualJOSN, _ := json.Marshal(finalScheduleEntity)
	assert.Equal(t, string(expectedJSON), string(actualJOSN))
}

func TestParsersSchedule_NewEntityFromModel(t *testing.T) {
	tenantID := "_"
	scheduleModel := testutils.FakeSchedule(tenantID)
	scheduleEntity := NewScheduleEntityFromModels(scheduleModel)
	finalScheduleModel := NewScheduleModelFromEntities(scheduleEntity)

	assert.Equal(t, finalScheduleModel.UUID, scheduleModel.UUID)
	assert.Equal(t, finalScheduleModel.TenantID, scheduleModel.TenantID)
	assert.Equal(t, finalScheduleModel.CreatedAt, scheduleModel.CreatedAt)
	assert.Equal(t, finalScheduleModel.Jobs[0].UUID, scheduleModel.Jobs[0].UUID)
	assert.Equal(t, finalScheduleModel.Jobs[0].Type, scheduleModel.Jobs[0].Type)
	assert.Equal(t, finalScheduleModel.Jobs[0].Labels, scheduleModel.Jobs[0].Labels)
}
