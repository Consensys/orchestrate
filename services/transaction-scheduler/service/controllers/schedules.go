package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/formatters"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
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
	router.Methods(http.MethodGet).Path("/schedules").HandlerFunc(c.get)
}

// @Summary Creates a new schedule
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 400
// @Failure 422
// @Failure 500
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

	scheduleEntity, err := c.ucs.CreateSchedule().Execute(ctx, &entities.Schedule{}, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	response := formatters.FormatScheduleResponse(scheduleEntity)
	_ = json.NewEncoder(rw).Encode(response)
}

// @Summary Fetch an schedule by its UUID
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 404
// @Failure 500
// @Router /schedules/{uuid} [get]
func (c *SchedulesController) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	uuid := mux.Vars(request)["uuid"]

	scheduleEntity, err := c.ucs.GetSchedule().Execute(ctx, uuid, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	response := formatters.FormatScheduleResponse(scheduleEntity)
	_ = json.NewEncoder(rw).Encode(response)
}

// @Summary Fetch a list of schedules
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 404
// @Failure 500
// @Router /schedules [get]
func (c *SchedulesController) get(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	scheduleEntities, err := c.ucs.GetSchedules().Execute(ctx, multitenancy.TenantIDFromContext(ctx))
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
