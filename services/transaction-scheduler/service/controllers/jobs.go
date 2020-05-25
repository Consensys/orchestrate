package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/formatters"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
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
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 400
// @Failure 500
// @Router /jobs [get]
func (c *JobsController) search(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	// @TODO Read filters from URL
	filters := make(map[string]string)

	jobRes, err := c.ucs.SearchJobs().Execute(ctx, filters, multitenancy.TenantIDFromContext(ctx))
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
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 400
// @Failure 500
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

	job := formatters.FormatJobCreateRequest(jobRequest)
	jobRes, err := c.ucs.CreateJob().Execute(ctx, job, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatJobResponse(jobRes))
}

// @Summary Fetch a job by its uuid
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 404
// @Failure 500
// @Router /jobs/{uuid} [get]
func (c *JobsController) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	uuid := mux.Vars(request)["uuid"]

	jobRes, err := c.ucs.GetJob().Execute(ctx, uuid, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatJobResponse(jobRes))
}

// @Summary Start a Job by its UUID
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 202
// @Failure 404
// @Failure 500
// @Router /jobs/{uuid}/start [put]
func (c *JobsController) start(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	jobUUID := mux.Vars(request)["uuid"]
	err := c.ucs.StartJob().Execute(ctx, jobUUID, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

// @Summary Update job by UUID
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 400
// @Failure 500
// @Router /jobs/{uuid} [path]
func (c *JobsController) update(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	jobRequest := &types.UpdateJobRequest{}
	err := jsonutils.UnmarshalBody(request.Body, jobRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	job := formatters.FormatJobUpdateRequest(jobRequest)
	job.UUID = mux.Vars(request)["uuid"]
	jobRes, err := c.ucs.UpdateJob().Execute(ctx, job, multitenancy.TenantIDFromContext(ctx))

	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatJobResponse(jobRes))
}
