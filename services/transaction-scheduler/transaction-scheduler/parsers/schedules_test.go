// +build unit

package parsers

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
)

func TestParsersSchedule_NewModelFromEntity(t *testing.T) {
	chainUUID := uuid.NewV4().String()
	tenantID := "_"
	scheduleEntity := testutils2.FakeScheduleEntity(chainUUID)
	scheduleModel := NewScheduleModelFromEntities(scheduleEntity, tenantID)
	finalScheduleEntity := NewScheduleEntityFromModels(scheduleModel)

	expectedJSON, _ := json.Marshal(scheduleEntity)
	actualJOSN, _ := json.Marshal(finalScheduleEntity)
	assert.Equal(t, string(expectedJSON), string(actualJOSN))
}

func TestParsersSchedule_NewEntityFromModel(t *testing.T) {
	tenantID := "_"
	scheduleModel := testutils.FakeSchedule(tenantID)
	scheduleEntity := NewScheduleEntityFromModels(scheduleModel)
	finalScheduleModel := NewScheduleModelFromEntities(scheduleEntity, tenantID)
	
	assert.Equal(t, finalScheduleModel.UUID, scheduleModel.UUID)
	assert.Equal(t, finalScheduleModel.TenantID, scheduleModel.TenantID)
	assert.Equal(t, finalScheduleModel.ChainUUID, scheduleModel.ChainUUID)
	assert.Equal(t, finalScheduleModel.CreatedAt, scheduleModel.CreatedAt)
	assert.Equal(t, finalScheduleModel.Jobs[0].UUID, scheduleModel.Jobs[0].UUID)
	assert.Equal(t, finalScheduleModel.Jobs[0].Type, scheduleModel.Jobs[0].Type)
	assert.Equal(t, finalScheduleModel.Jobs[0].Labels, scheduleModel.Jobs[0].Labels)
}