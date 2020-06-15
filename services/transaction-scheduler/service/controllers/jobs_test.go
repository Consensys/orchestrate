// +build unit

package controllers

import (
	"bytes"
	"context"
	"fmt"
	types2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	testutils3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs/mocks"
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

func (s jobsCtrlTestSuite) CreateJob() jobs.CreateJobUseCase {
	return s.createJobUC
}

func (s jobsCtrlTestSuite) GetJob() jobs.GetJobUseCase {
	return s.getJobUC
}

func (s jobsCtrlTestSuite) StartJob() jobs.StartJobUseCase {
	return s.startJobUC
}

func (s jobsCtrlTestSuite) UpdateJob() jobs.UpdateJobUseCase {
	return s.updateJobUC
}

func (s jobsCtrlTestSuite) SearchJobs() jobs.SearchJobsUseCase {
	return s.searchJobUC
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
		jobRequest := testutils.FakeCreateJobRequest()
		jobEntityRes := testutils3.FakeJob()
		requestBytes, _ := json.Marshal(jobRequest)
		rw := httptest.NewRecorder()

		httpRequest := httptest.
			NewRequest(http.MethodPost, "/jobs", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.createJobUC.EXPECT().
			Execute(gomock.Any(), gomock.Any(), s.tenantID).
			Return(jobEntityRes, nil).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatJobResponse(jobEntityRes)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with Bad request if invalid format", func(t *testing.T) {
		jobRequest := testutils.FakeCreateJobRequest()
		jobRequest.ScheduleUUID = ""
		requestBytes, _ := json.Marshal(jobRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/jobs", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails with InvalidParameterError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		jobRequest := testutils.FakeCreateJobRequest()
		requestBytes, _ := json.Marshal(jobRequest)
		httpRequest := httptest.
			NewRequest(http.MethodPost, "/jobs", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.createJobUC.EXPECT().
			Execute(gomock.Any(), gomock.Any(), s.tenantID).
			Return(nil, errors.InvalidParameterError("error")).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}

func (s *jobsCtrlTestSuite) TestJobsController_GetOne() {
	s.T().Run("should execute get one job request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodGet, "/jobs/jobUUID", nil).
			WithContext(s.ctx)
		jobEntityRes := testutils3.FakeJob()

		s.getJobUC.EXPECT().
			Execute(gomock.Any(), "jobUUID", s.tenantID).
			Return(jobEntityRes, nil).Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatJobResponse(jobEntityRes)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 404 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodGet, "/schedules/jobUUID", bytes.NewReader(nil)).
			WithContext(s.ctx)
		s.getJobUC.EXPECT().
			Execute(gomock.Any(), "jobUUID", s.tenantID).
			Return(nil, errors.NotFoundError("error")).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *jobsCtrlTestSuite) TestJobsController_Search() {
	s.T().Run("should execute search jobs successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		filters := &entities.JobFilters{}
		httpRequest := httptest.
			NewRequest(http.MethodGet, "/jobs", nil).
			WithContext(s.ctx)
		jobEntities := []*types2.Job{testutils3.FakeJob()}

		s.searchJobUC.EXPECT().
			Execute(gomock.Any(), filters, s.tenantID).
			Return(jobEntities, nil).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		response := []*types.JobResponse{formatters.FormatJobResponse(jobEntities[0])}
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should execute search jobs by tx_hashes successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		filters := &entities.JobFilters{
			TxHashes: []string{
				common.HexToHash("0x1").String(),
				common.HexToHash("0x2").String(),
			},
		}
		url := fmt.Sprintf("/jobs?tx_hashes=%s", strings.Join([]string{
			common.HexToHash("0x1").String(),
			common.HexToHash("0x2").String(),
		}, ","))

		httpRequest := httptest.
			NewRequest(http.MethodGet, url, nil).
			WithContext(s.ctx)
		jobEntities := []*types2.Job{testutils3.FakeJob()}

		s.searchJobUC.EXPECT().
			Execute(gomock.Any(), filters, s.tenantID).
			Return(jobEntities, nil).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		response := []*types.JobResponse{formatters.FormatJobResponse(jobEntities[0])}
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails on invalid tx hashes as input", func(t *testing.T) {
		rw := httptest.NewRecorder()
		url := fmt.Sprintf("/jobs?tx_hashes=%s", strings.Join([]string{
			"InvalidHash",
			common.HexToHash("0x2").String(),
		}, ","))

		httpRequest := httptest.
			NewRequest(http.MethodGet, url, bytes.NewReader(nil)).
			WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		filters := &entities.JobFilters{}
		httpRequest := httptest.
			NewRequest(http.MethodGet, "/jobs", bytes.NewReader(nil)).
			WithContext(s.ctx)

		s.searchJobUC.EXPECT().
			Execute(gomock.Any(), filters, s.tenantID).
			Return(nil, errors.InvalidParameterError("error")).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}

func (s *jobsCtrlTestSuite) TestJobsController_Start() {
	s.T().Run("should execute start a job request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPut, "/jobs/jobUUID/start", nil).
			WithContext(s.ctx)

		s.startJobUC.EXPECT().
			Execute(gomock.Any(), "jobUUID", s.tenantID).
			Return(nil).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, http.StatusAccepted, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 404 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPut, "/jobs/jobUUID/start", bytes.NewReader(nil)).
			WithContext(s.ctx)

		s.startJobUC.EXPECT().
			Execute(gomock.Any(), "jobUUID", s.tenantID).
			Return(errors.NotFoundError("error")).
			Times(1)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *jobsCtrlTestSuite) TestJobsController_Update() {
	s.T().Run("should execute update a job request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		jobRequest := testutils.FakeJobUpdateRequest()
		jobEntityRes := testutils3.FakeJob()

		requestBytes, _ := json.Marshal(jobRequest)
		httpRequest := httptest.
			NewRequest(http.MethodPatch, "/jobs/"+jobEntityRes.UUID, bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		jobEntityReq := formatters.FormatJobUpdateRequest(jobRequest)
		jobEntityReq.UUID = jobEntityRes.UUID
		s.updateJobUC.EXPECT().
			Execute(gomock.Any(), gomock.Any(), jobRequest.Status, s.tenantID).
			Return(jobEntityRes, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatJobResponse(jobEntityRes)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 422 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPatch, "/jobs/jobUUID", bytes.NewReader(nil)).
			WithContext(s.ctx)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})
}
