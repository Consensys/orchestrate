// +build unit

package controllers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/formatters"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules/mocks"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
)

type schedulesCtrlTestSuite struct {
	suite.Suite
	createScheduleUC *mocks.MockCreateScheduleUseCase
	getScheduleUC    *mocks.MockGetScheduleUseCase
	getSchedulesUC   *mocks.MockGetSchedulesUseCase
	ctx              context.Context
	tenantID         string
	router           *mux.Router
}

func (s *schedulesCtrlTestSuite) CreateSchedule() schedules.CreateScheduleUseCase {
	return s.createScheduleUC
}

func (s *schedulesCtrlTestSuite) GetSchedule() schedules.GetScheduleUseCase {
	return s.getScheduleUC
}

func (s *schedulesCtrlTestSuite) GetSchedules() schedules.GetSchedulesUseCase {
	return s.getSchedulesUC
}

var _ schedules.UseCases = &schedulesCtrlTestSuite{}

func TestSchedulesController(t *testing.T) {
	s := new(schedulesCtrlTestSuite)
	suite.Run(t, s)
}

func (s *schedulesCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.tenantID = "tenantID"
	s.createScheduleUC = mocks.NewMockCreateScheduleUseCase(ctrl)
	s.getScheduleUC = mocks.NewMockGetScheduleUseCase(ctrl)
	s.getSchedulesUC = mocks.NewMockGetSchedulesUseCase(ctrl)
	s.ctx = context.WithValue(context.Background(), multitenancy.TenantIDKey, s.tenantID)
	s.router = mux.NewRouter()

	controller := NewSchedulesController(s)
	controller.Append(s.router)
}

func (s *schedulesCtrlTestSuite) TestScheduleController_Create() {
	scheduleRequest := testutils.FakeCreateScheduleRequest()
	scheduleEntity := formatters.FormatScheduleCreateRequest(scheduleRequest)
	requestBytes, _ := json.Marshal(scheduleRequest)

	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/schedules", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)
		
		scheduleEntityResp := testutils2.FakeScheduleEntity(scheduleRequest.ChainUUID)

		s.createScheduleUC.EXPECT().
			Execute(gomock.Any(), scheduleEntity, s.tenantID).
			Return(scheduleEntityResp, nil).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatScheduleResponse(scheduleEntityResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		scheduleRequest := testutils.FakeCreateScheduleRequest()
		scheduleRequest.ChainUUID = ""
		requestBytes, _ := json.Marshal(scheduleRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/schedules", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails with InvalidParameterError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/schedules", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.createScheduleUC.EXPECT().
			Execute(gomock.Any(), scheduleEntity, s.tenantID).
			Return(nil, errors.InvalidParameterError("error")).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}


func (s *schedulesCtrlTestSuite) TestScheduleController_GetOne() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/schedules/scheduleUUID", nil).WithContext(s.ctx)
		scheduleEntityResp := testutils2.FakeScheduleEntity("ChainUUID")

		s.getScheduleUC.EXPECT().
			Execute(gomock.Any(), "scheduleUUID", s.tenantID).
			Return(scheduleEntityResp, nil).Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatScheduleResponse(scheduleEntityResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 404 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/schedules/scheduleUUID", bytes.NewReader(nil)).
			WithContext(s.ctx)

		s.getScheduleUC.EXPECT().
			Execute(gomock.Any(), "scheduleUUID", s.tenantID).
			Return(nil, errors.NotFoundError("error")).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *schedulesCtrlTestSuite) TestScheduleController_GetAll() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/schedules", nil).WithContext(s.ctx)
		schedulesEntities := []*entities.Schedule{testutils2.FakeScheduleEntity("chainUUID")}

		s.getSchedulesUC.EXPECT().
			Execute(gomock.Any(), s.tenantID).
			Return(schedulesEntities, nil).
			Times(1)
		s.router.ServeHTTP(rw, httpRequest)

		response := []*types.ScheduleResponse{formatters.FormatScheduleResponse(schedulesEntities[0])}
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 404 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodGet, "/schedules/scheduleUUID", bytes.NewReader(nil)).
			WithContext(s.ctx)
		s.getScheduleUC.EXPECT().
			Execute(gomock.Any(), "scheduleUUID", s.tenantID).
			Return(nil, errors.NotFoundError("error")).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}
