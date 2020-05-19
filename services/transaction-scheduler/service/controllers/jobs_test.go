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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
)

type jobsCtrlTestSuite struct {
	suite.Suite
	createJobUC *mocks.MockCreateJobUseCase
	getJobUC    *mocks.MockGetJobUseCase
	startJobUC  *mocks.MockStartJobUseCase
	updateJobUC *mocks.MockUpdateJobUseCase
	searchJobUC *mocks.MockSearchJobsUseCase
	ctx         context.Context
	tenantID    string
	router      *mux.Router
}

var _ jobs.UseCases = &jobsCtrlTestSuite{}

func (t jobsCtrlTestSuite) CreateJob() jobs.CreateJobUseCase {
	return t.createJobUC
}

func (t jobsCtrlTestSuite) GetJob() jobs.GetJobUseCase {
	return t.getJobUC
}

func (t jobsCtrlTestSuite) StartJob() jobs.StartJobUseCase {
	return t.startJobUC
}

func (t jobsCtrlTestSuite) UpdateJob() jobs.UpdateJobUseCase {
	return t.updateJobUC
}

func (t jobsCtrlTestSuite) SearchJobs() jobs.SearchJobsUseCase {
	return t.searchJobUC
}

func TestJobsController(t *testing.T) {
	s := new(jobsCtrlTestSuite)
	suite.Run(t, s)
}

func (s *jobsCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.tenantID = "tenantID"
	s.createJobUC = mocks.NewMockCreateJobUseCase(ctrl)
	s.getJobUC = mocks.NewMockGetJobUseCase(ctrl)
	s.startJobUC = mocks.NewMockStartJobUseCase(ctrl)
	s.updateJobUC = mocks.NewMockUpdateJobUseCase(ctrl)
	s.searchJobUC = mocks.NewMockSearchJobsUseCase(ctrl)
	s.ctx = context.WithValue(context.Background(), multitenancy.TenantIDKey, s.tenantID)
	s.router = mux.NewRouter()

	controller := NewJobsController(s)
	controller.Append(s.router)
}

func (s *jobsCtrlTestSuite) TestJobsController_Create() {
	s.T().Run("should execute create job request successfully", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		requestBytes, _ := json.Marshal(jobRequest)
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewReader(requestBytes)).WithContext(s.ctx)
		jobResponse := testutils.FakeJobResponse()

		s.createJobUC.EXPECT().Execute(gomock.Any(), jobRequest, s.tenantID).Return(jobResponse, nil).Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := utils.ObjectToJSON(jobResponse)
		assert.Equal(t, expectedBody+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		jobRequest := testutils.FakeJobRequest()
		jobRequest.ScheduleUUID = ""
		requestBytes, _ := json.Marshal(jobRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails with InvalidParameterError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		jobRequest := testutils.FakeJobRequest()
		requestBytes, _ := json.Marshal(jobRequest)
		httpRequest := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)
		s.createJobUC.EXPECT().Execute(gomock.Any(), jobRequest, s.tenantID).
			Return(nil, errors.InvalidParameterError("error")).Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}

func (s *jobsCtrlTestSuite) TestJobsController_GetOne() {
	s.T().Run("should execute get one job request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/jobs/jobUUID", nil).WithContext(s.ctx)
		jobResponse := testutils.FakeJobResponse()

		s.getJobUC.EXPECT().Execute(gomock.Any(), "jobUUID", s.tenantID).Return(jobResponse, nil).Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := utils.ObjectToJSON(jobResponse)
		assert.Equal(t, expectedBody+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 404 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/schedules/jobUUID", bytes.NewReader(nil)).
			WithContext(s.ctx)
		s.getJobUC.EXPECT().Execute(gomock.Any(), "jobUUID", s.tenantID).
			Return(nil, errors.NotFoundError("error")).Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *jobsCtrlTestSuite) TestJobsController_Search() {
	s.T().Run("should execute search jobs successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/jobs", nil).WithContext(s.ctx)
		jobsResponse := []*types.JobResponse{testutils.FakeJobResponse()}

		s.searchJobUC.EXPECT().Execute(gomock.Any(), map[string]string{}, s.tenantID).
			Return(jobsResponse, nil).Times(1)
		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := utils.ObjectToJSON(jobsResponse)
		assert.Equal(t, expectedBody+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/jobs", bytes.NewReader(nil)).WithContext(s.ctx)
		s.searchJobUC.EXPECT().Execute(gomock.Any(), map[string]string{}, s.tenantID).
			Return(nil, errors.InvalidParameterError("error")).Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}

func (s *jobsCtrlTestSuite) TestJobsController_Start() {
	s.T().Run("should execute start a job request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPut, "/jobs/jobUUID/start", nil).WithContext(s.ctx)

		s.startJobUC.EXPECT().Execute(gomock.Any(), "jobUUID", s.tenantID).
			Return(nil).Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, http.StatusAccepted, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 404 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPut, "/jobs/jobUUID/start", bytes.NewReader(nil)).
			WithContext(s.ctx)
		s.startJobUC.EXPECT().Execute(gomock.Any(), "jobUUID", s.tenantID).
			Return(errors.NotFoundError("error")).Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *jobsCtrlTestSuite) TestJobsController_Update() {
	s.T().Run("should execute update a job request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		jobRequest := testutils.FakeJobUpdateRequest()
		requestBytes, _ := json.Marshal(jobRequest)
		httpRequest := httptest.NewRequest(http.MethodPatch, "/jobs/jobUUID", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)
		jobResponse := testutils.FakeJobResponse()

		s.updateJobUC.EXPECT().Execute(gomock.Any(), "jobUUID", jobRequest, s.tenantID).
			Return(jobResponse, nil).Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := utils.ObjectToJSON(jobResponse)
		assert.Equal(t, expectedBody+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPatch, "/jobs/jobUUID", bytes.NewReader(nil)).
			WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})
}
