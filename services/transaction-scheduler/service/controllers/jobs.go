package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/formatters"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
)

type JobsController struct {
	ucs jobs.UseCases
}

func NewJobsController(useCases jobs.UseCases) *JobsController {
	return &JobsController{
		ucs: useCases,
	}
}

// Add routes to router
func (c *JobsController) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/jobs").HandlerFunc(c.search)
	router.Methods(http.MethodPost).Path("/jobs").HandlerFunc(c.create)
	router.Methods(http.MethodGet).Path("/jobs/{uuid}").HandlerFunc(c.getOne)
	router.Methods(http.MethodPatch).Path("/jobs/{uuid}").HandlerFunc(c.update)
	router.Methods(http.MethodPut).Path("/jobs/{uuid}/start").HandlerFunc(c.start)
}

// @Summary Search jobs by provided filters
// @Description Get a list of filtered jobs
// @Tags Jobs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param tx_hashes query []string false "List of transaction hashes" collectionFormat(csv)
// @Param chain_uuid query string false "Chain UUID"
// @Success 200 {object} types.JobResponse{annotations=types.Annotations,transaction=types.ETHTransaction,logs=[]types.Log} "List of Jobs found"
// @Failure 400 {string} httputil.ErrorResponse "Invalid filter in the request"
// @Failure 500 {string} httputil.ErrorResponse "Internal server error"
// @Router /jobs [get]
func (c *JobsController) search(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	filters, err := formatters.FormatJobFilterRequest(request)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	jobRes, err := c.ucs.SearchJobs().Execute(ctx, filters, multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	var response []*types.JobResponse
	for _, jb := range jobRes {
		response = append(response, formatters.FormatJobResponse(jb))
	}

	_ = json.NewEncoder(rw).Encode(response)
}

// @Summary Creates a new Job
// @Description Creates a new job as part of an already created schedule
// @Tags Jobs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body types.CreateJobRequest{annotations=types.Annotations,transaction=types.ETHTransaction} true "Job creation request"
// @Success 200 {object} types.JobResponse{annotations=types.Annotations,transaction=types.ETHTransaction,logs=[]types.Log} "Created Job"
// @Failure 400 {string} httputil.ErrorResponse "Invalid request"
// @Failure 422 {string} httputil.ErrorResponse "Unprocessable parameters were sent"
// @Failure 500 {string} httputil.ErrorResponse "Internal server error"
// @Router /jobs [post]
func (c *JobsController) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	jobRequest := &types.CreateJobRequest{}
	err := jsonutils.UnmarshalBody(request.Body, jobRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err = jobRequest.Annotations.Validate(); err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	job := formatters.FormatJobCreateRequest(jobRequest)
	jobRes, err := c.ucs.CreateJob().Execute(ctx, job, multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatJobResponse(jobRes))
}

// @Summary Fetch a job by uuid
// @Description Fetch a single job by uuid
// @Tags Jobs
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "UUID of the job"
// @Success 200 {object} types.JobResponse{annotations=types.Annotations,transaction=types.ETHTransaction,logs=[]types.Log} "Job found"
// @Failure 404 {string} httputil.ErrorResponse "Job not found"
// @Failure 500 {string} httputil.ErrorResponse "Internal server error"
// @Router /jobs/{uuid} [get]
func (c *JobsController) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	uuid := mux.Vars(request)["uuid"]

	jobRes, err := c.ucs.GetJob().Execute(ctx, uuid, multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatJobResponse(jobRes))
}

// @Summary Start a Job by UUID
// @Description Starts a specific job by UUID, effectively executing the transaction asynchronously
// @Tags Jobs
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "UUID of the job"
// @Success 202
// @Failure 404 {string} httputil.ErrorResponse "Job not found"
// @Failure 500 {string} httputil.ErrorResponse "Internal server error"
// @Router /jobs/{uuid}/start [put]
func (c *JobsController) start(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	jobUUID := mux.Vars(request)["uuid"]
	err := c.ucs.StartJob().Execute(ctx, jobUUID, multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

// @Summary Update job by UUID
// @Description Update a specific job by UUID
// @Description WARNING: Reserved for advanced users. Orchestrate does not recommend using this endpoint.
// @Tags Jobs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body types.UpdateJobRequest{annotations=types.Annotations,transaction=types.ETHTransaction} true "Job update request"
// @Success 200 {object} types.JobResponse{annotations=types.Annotations,transaction=types.ETHTransaction,logs=[]types.Log} "Job found"
// @Failure 400 {string} httputil.ErrorResponse "Invalid request"
// @Failure 404 {string} httputil.ErrorResponse "Job not found"
// @Failure 409 {string} httputil.ErrorResponse "Job in invalid state for the given status update"
// @Failure 500 {string} httputil.ErrorResponse "Internal server error"
// @Router /jobs/{uuid} [patch]
func (c *JobsController) update(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	jobRequest := &types.UpdateJobRequest{}
	err := jsonutils.UnmarshalBody(request.Body, jobRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err = jobRequest.Annotations.Validate(); err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	job := formatters.FormatJobUpdateRequest(jobRequest)
	job.UUID = mux.Vars(request)["uuid"]
	jobRes, err := c.ucs.UpdateJob().Execute(ctx, job, jobRequest.Status, jobRequest.Message,
		multitenancy.AllowedTenantsFromContext(ctx))

	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatJobResponse(jobRes))
}
