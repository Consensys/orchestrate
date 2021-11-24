// +build unit

package controllers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	txschedulertypes "github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"

	"github.com/consensys/orchestrate/pkg/encoding/json"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/services/api/business/use-cases/mocks"
	"github.com/consensys/orchestrate/services/api/service/formatters"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type schedulesCtrlTestSuite struct {
	suite.Suite
	createScheduleUC  *mocks.MockCreateScheduleUseCase
	getScheduleUC     *mocks.MockGetScheduleUseCase
	searchSchedulesUC *mocks.MockSearchSchedulesUseCase
	ctx               context.Context
	userInfo          *multitenancy.UserInfo
	router            *mux.Router
}

func (s *schedulesCtrlTestSuite) CreateSchedule() usecases.CreateScheduleUseCase {
	return s.createScheduleUC
}

func (s *schedulesCtrlTestSuite) GetSchedule() usecases.GetScheduleUseCase {
	return s.getScheduleUC
}

func (s *schedulesCtrlTestSuite) SearchSchedules() usecases.SearchSchedulesUseCase {
	return s.searchSchedulesUC
}

var _ usecases.ScheduleUseCases = &schedulesCtrlTestSuite{}

func TestSchedulesController(t *testing.T) {
	s := new(schedulesCtrlTestSuite)
	suite.Run(t, s)
}

func (s *schedulesCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.createScheduleUC = mocks.NewMockCreateScheduleUseCase(ctrl)
	s.getScheduleUC = mocks.NewMockGetScheduleUseCase(ctrl)
	s.searchSchedulesUC = mocks.NewMockSearchSchedulesUseCase(ctrl)
	s.userInfo = multitenancy.NewUserInfo("tenantOne", "username")
	s.ctx = multitenancy.WithUserInfo(context.Background(), s.userInfo)
	s.router = mux.NewRouter()

	controller := NewSchedulesController(s)
	controller.Append(s.router)
}

func (s *schedulesCtrlTestSuite) TestScheduleController_Create() {
	scheduleRequest := testutils.FakeCreateScheduleRequest()
	requestBytes, _ := json.Marshal(scheduleRequest)

	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/schedules", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		scheduleEntityResp := testutils.FakeSchedule()

		s.createScheduleUC.EXPECT().
			Execute(gomock.Any(), &entities.Schedule{}, s.userInfo).
			Return(scheduleEntityResp, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatScheduleResponse(scheduleEntityResp)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails with InvalidParameterError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/schedules", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.createScheduleUC.EXPECT().
			Execute(gomock.Any(), &entities.Schedule{}, s.userInfo).
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
		scheduleEntityResp := testutils.FakeSchedule()

		s.getScheduleUC.EXPECT().
			Execute(gomock.Any(), "scheduleUUID", s.userInfo).
			Return(scheduleEntityResp, nil)

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
			Execute(gomock.Any(), "scheduleUUID", s.userInfo).
			Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *schedulesCtrlTestSuite) TestScheduleController_GetAll() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/schedules", nil).WithContext(s.ctx)
		schedulesEntities := []*entities.Schedule{testutils.FakeSchedule()}

		s.searchSchedulesUC.EXPECT().Execute(gomock.Any(), s.userInfo).Return(schedulesEntities, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := []*txschedulertypes.ScheduleResponse{formatters.FormatScheduleResponse(schedulesEntities[0])}
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
			Execute(gomock.Any(), "scheduleUUID", s.userInfo).
			Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}
