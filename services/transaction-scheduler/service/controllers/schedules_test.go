// +build unit

package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
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
	scheduleRequest := testutils.FakeScheduleRequest()
	requestBytes, _ := json.Marshal(scheduleRequest)

	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/schedules", bytes.NewReader(requestBytes)).WithContext(s.ctx)
		scheduleResponse := testutils.FakeScheduleResponse()

		s.createScheduleUC.EXPECT().Execute(gomock.Any(), scheduleRequest, s.tenantID).Return(scheduleResponse, nil).Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := utils.ObjectToJSON(scheduleResponse)
		assert.Equal(t, expectedBody+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		scheduleRequest := testutils.FakeScheduleRequest()
		scheduleRequest.ChainUUID = ""
		requestBytes, _ := json.Marshal(scheduleRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/schedules", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails with InvalidParameterError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/schedules", bytes.NewReader(requestBytes)).WithContext(s.ctx)
		s.createScheduleUC.EXPECT().Execute(gomock.Any(), scheduleRequest, s.tenantID).Return(nil, errors.InvalidParameterError("error")).Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}

func (s *schedulesCtrlTestSuite) TestScheduleController_GetOne() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/schedules/scheduleUUID", nil).WithContext(s.ctx)
		scheduleResponse := testutils.FakeScheduleResponse()

		s.getScheduleUC.EXPECT().Execute(gomock.Any(), "scheduleUUID", s.tenantID).Return(scheduleResponse, nil).Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := utils.ObjectToJSON(scheduleResponse)
		assert.Equal(t, expectedBody+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 404 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/schedules/scheduleUUID", bytes.NewReader(nil)).WithContext(s.ctx)
		s.getScheduleUC.EXPECT().Execute(gomock.Any(), "scheduleUUID", s.tenantID).
			Return(nil, errors.NotFoundError("error")).Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *schedulesCtrlTestSuite) TestScheduleController_GetAll() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/schedules", nil).WithContext(s.ctx)
		schedulesResponse := []*types.ScheduleResponse{testutils.FakeScheduleResponse()}

		s.getSchedulesUC.EXPECT().Execute(gomock.Any(), s.tenantID).
			Return(schedulesResponse, nil).Times(1)
		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := utils.ObjectToJSON(schedulesResponse)
		assert.Equal(t, expectedBody+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 404 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/schedules/scheduleUUID", bytes.NewReader(nil)).WithContext(s.ctx)
		s.getScheduleUC.EXPECT().Execute(gomock.Any(), "scheduleUUID", s.tenantID).
			Return(nil, errors.NotFoundError("error")).Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}
