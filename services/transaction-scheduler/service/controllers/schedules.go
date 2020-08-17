package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx-scheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/formatters"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules"
)

type SchedulesController struct {
	ucs schedules.UseCases
}

func NewSchedulesController(useCases schedules.UseCases) *SchedulesController {
	return &SchedulesController{
		ucs: useCases,
	}
}

// Add routes to router
func (c *SchedulesController) Append(router *mux.Router) {
	router.Methods(http.MethodPost).Path("/schedules").HandlerFunc(c.create)
	router.Methods(http.MethodGet).Path("/schedules/{uuid}").HandlerFunc(c.getOne)
	router.Methods(http.MethodGet).Path("/schedules").HandlerFunc(c.getAll)
}

// @Summary Creates a new Schedule
// @Description Creates a new schedule
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body types.CreateScheduleRequest true "Schedule creation request"
// @Success 200 {object} types.ScheduleResponse{jobs=[]types.JobResponse} "Created schedule"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 422 {object} httputil.ErrorResponse "Unprocessable parameters were sent"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /schedules [post]
func (c *SchedulesController) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	scheduleRequest := &types.CreateScheduleRequest{}
	err := jsonutils.UnmarshalBody(request.Body, scheduleRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	scheduleEntity, err := c.ucs.CreateSchedule().Execute(ctx, &entities.Schedule{TenantID: multitenancy.TenantIDFromContext(ctx)})
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	response := formatters.FormatScheduleResponse(scheduleEntity)
	_ = json.NewEncoder(rw).Encode(response)
}

// @Summary Fetch a schedule by uuid
// @Description Fetch a single schedule by uuid
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "UUID of the schedule"
// @Success 200 {object} types.ScheduleResponse "Schedule found"
// @Failure 404 {object} httputil.ErrorResponse "Schedule not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /schedules/{uuid} [get]
func (c *SchedulesController) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	uuid := mux.Vars(request)["uuid"]

	scheduleEntity, err := c.ucs.GetSchedule().Execute(ctx, uuid, multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	response := formatters.FormatScheduleResponse(scheduleEntity)
	_ = json.NewEncoder(rw).Encode(response)
}

// @Summary Get all schedules
// @Description Get all schedules
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200 {array} types.ScheduleResponse "List of schedules found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /schedules [get]
func (c *SchedulesController) getAll(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	scheduleEntities, err := c.ucs.SearchSchedules().Execute(ctx, multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	var response []*types.ScheduleResponse
	for _, scheduleEntity := range scheduleEntities {
		response = append(response, formatters.FormatScheduleResponse(scheduleEntity))
	}

	_ = json.NewEncoder(rw).Encode(response)
}
