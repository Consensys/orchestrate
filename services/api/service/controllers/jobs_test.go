// +build unit

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	txschedulertypes "github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/business/use-cases/mocks"

	"encoding/json"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/formatters"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type jobsCtrlTestSuite struct {
	suite.Suite
	createJobUC    *mocks.MockCreateJobUseCase
	getJobUC       *mocks.MockGetJobUseCase
	startJobUC     *mocks.MockStartJobUseCase
	resentJobTxUC  *mocks.MockResendJobTxUseCase
	startNextJobUC *mocks.MockStartNextJobUseCase
	updateJobUC    *mocks.MockUpdateJobUseCase
	searchJobUC    *mocks.MockSearchJobsUseCase
	ctx            context.Context
	userInfo       *multitenancy.UserInfo
	router         *mux.Router
}

var _ usecases.JobUseCases = &jobsCtrlTestSuite{}

func (s jobsCtrlTestSuite) CreateJob() usecases.CreateJobUseCase {
	return s.createJobUC
}

func (s jobsCtrlTestSuite) GetJob() usecases.GetJobUseCase {
	return s.getJobUC
}

func (s jobsCtrlTestSuite) StartJob() usecases.StartJobUseCase {
	return s.startJobUC
}

func (s jobsCtrlTestSuite) ResendJobTx() usecases.ResendJobTxUseCase {
	return s.resentJobTxUC
}

func (s jobsCtrlTestSuite) StartNextJob() usecases.StartNextJobUseCase {
	return s.startNextJobUC
}

func (s jobsCtrlTestSuite) UpdateJob() usecases.UpdateJobUseCase {
	return s.updateJobUC
}

func (s jobsCtrlTestSuite) SearchJobs() usecases.SearchJobsUseCase {
	return s.searchJobUC
}

func TestJobsController(t *testing.T) {
	s := new(jobsCtrlTestSuite)
	suite.Run(t, s)
}

func (s *jobsCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.createJobUC = mocks.NewMockCreateJobUseCase(ctrl)
	s.getJobUC = mocks.NewMockGetJobUseCase(ctrl)
	s.startJobUC = mocks.NewMockStartJobUseCase(ctrl)
	s.updateJobUC = mocks.NewMockUpdateJobUseCase(ctrl)
	s.searchJobUC = mocks.NewMockSearchJobsUseCase(ctrl)
	s.resentJobTxUC = mocks.NewMockResendJobTxUseCase(ctrl)
	s.userInfo = multitenancy.NewUserInfo("tenantOne", "username")
	s.ctx = multitenancy.WithUserInfo(context.Background(), s.userInfo)
	s.router = mux.NewRouter()

	controller := NewJobsController(s)
	controller.Append(s.router)
}

func (s *jobsCtrlTestSuite) TestJobsController_Create() {
	s.T().Run("should execute create job request successfully", func(t *testing.T) {
		jobRequest := testutils.FakeCreateJobRequest()
		jobRequest.Annotations = txschedulertypes.Annotations{
			OneTimeKey: true,
		}
		jobEntityRes := testutils.FakeJob()
		requestBytes, _ := json.Marshal(jobRequest)
		rw := httptest.NewRecorder()

		httpRequest := httptest.
			NewRequest(http.MethodPost, "/jobs", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.createJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), s.userInfo).Return(jobEntityRes, nil)

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

	s.T().Run("should fail with Bad request if invalid format (retry)", func(t *testing.T) {
		jobRequest := testutils.FakeCreateJobRequest()
		jobRequest.Annotations = txschedulertypes.Annotations{
			GasPricePolicy: txschedulertypes.GasPriceParams{
				RetryPolicy: txschedulertypes.RetryParams{
					Interval:  "invalid",
					Increment: 1.1,
					Limit:     1.4,
				},
			},
		}
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

		s.createJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), s.userInfo).Return(nil, errors.InvalidParameterError("error"))

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
		jobEntityRes := testutils.FakeJob()

		s.getJobUC.EXPECT().Execute(gomock.Any(), "jobUUID", s.userInfo).Return(jobEntityRes, nil)

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
		s.getJobUC.EXPECT().Execute(gomock.Any(), "jobUUID", s.userInfo).Return(nil, errors.NotFoundError("error"))
		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *jobsCtrlTestSuite) TestJobsController_Search() {
	s.T().Run("should execute search jobs successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		filters := &entities.JobFilters{}
		httpRequest := httptest.NewRequest(http.MethodGet, "/jobs", nil).WithContext(s.ctx)
		jobEntities := []*entities.Job{testutils.FakeJob()}

		s.searchJobUC.EXPECT().Execute(gomock.Any(), filters, s.userInfo).Return(jobEntities, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := []*txschedulertypes.JobResponse{formatters.FormatJobResponse(jobEntities[0])}
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should execute search jobs by tx_hashes successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		filters := &entities.JobFilters{TxHashes: []string{common.HexToHash("0x1").String(), common.HexToHash("0x2").String()}}
		url := fmt.Sprintf("/jobs?tx_hashes=%s", strings.Join([]string{
			common.HexToHash("0x1").String(),
			common.HexToHash("0x2").String(),
		}, ","))

		httpRequest := httptest.
			NewRequest(http.MethodGet, url, nil).
			WithContext(s.ctx)
		jobEntities := []*entities.Job{testutils.FakeJob()}

		s.searchJobUC.EXPECT().Execute(gomock.Any(), filters, s.userInfo).Return(jobEntities, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := []*txschedulertypes.JobResponse{formatters.FormatJobResponse(jobEntities[0])}
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

		s.searchJobUC.EXPECT().Execute(gomock.Any(), filters, s.userInfo).Return(nil, errors.InvalidParameterError("error"))

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

		s.startJobUC.EXPECT().Execute(gomock.Any(), "jobUUID", s.userInfo).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, http.StatusAccepted, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 404 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPut, "/jobs/jobUUID/start", bytes.NewReader(nil)).
			WithContext(s.ctx)

		s.startJobUC.EXPECT().Execute(gomock.Any(), "jobUUID", s.userInfo).Return(errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *jobsCtrlTestSuite) TestJobsController_ResendTxJob() {
	s.T().Run("should execute resend job transaction request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPut, "/jobs/jobUUID/resend", nil).
			WithContext(s.ctx)

		s.resentJobTxUC.EXPECT().Execute(gomock.Any(), "jobUUID", s.userInfo).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, http.StatusAccepted, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 404 if use case fails with NotFoundError", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPut, "/jobs/jobUUID/resend", bytes.NewReader(nil)).
			WithContext(s.ctx)

		s.resentJobTxUC.EXPECT().Execute(gomock.Any(), "jobUUID", s.userInfo).Return(errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *jobsCtrlTestSuite) TestJobsController_Update() {
	s.T().Run("should execute update a job request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		jobRequest := testutils.FakeJobUpdateRequest()
		jobRequest.Annotations = &txschedulertypes.Annotations{OneTimeKey: true}
		jobEntityRes := testutils.FakeJob()

		requestBytes, _ := json.Marshal(jobRequest)
		httpRequest := httptest.
			NewRequest(http.MethodPatch, "/jobs/"+jobEntityRes.UUID, bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		jobEntityReq := formatters.FormatJobUpdateRequest(jobRequest)
		jobEntityReq.UUID = jobEntityRes.UUID
		s.updateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), jobRequest.Status, "", s.userInfo).Return(jobEntityRes, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatJobResponse(jobEntityRes)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with 409 if use case fails with InvalidStateError", func(t *testing.T) {
		jobRequest := testutils.FakeJobUpdateRequest()
		requestBytes, _ := json.Marshal(jobRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.
			NewRequest(http.MethodPatch, "/jobs/jobUUID", bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.updateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any(), "", s.userInfo).Return(nil, errors.InvalidStateError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusConflict, rw.Code)
	})
}
