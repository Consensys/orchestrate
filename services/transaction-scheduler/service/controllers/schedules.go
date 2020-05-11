package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

type SchedulesController struct {
	createScheduleUseCase schedules.CreateScheduleUseCase
	getScheduleUseCase    schedules.GetScheduleUseCase
}

func NewSchedulesController(createScheduleUseCase schedules.CreateScheduleUseCase, getScheduleUseCase schedules.GetScheduleUseCase) *SchedulesController {
	return &SchedulesController{
		createScheduleUseCase: createScheduleUseCase,
		getScheduleUseCase:    getScheduleUseCase,
	}
}

// Add routes to router
func (c *SchedulesController) Append(router *mux.Router) {
	router.Methods(http.MethodPost).Path("/schedules").HandlerFunc(c.Create)
	router.Methods(http.MethodGet).Path("/schedules/{uuid}").HandlerFunc(c.Get)
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
func (c *SchedulesController) Create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	scheduleRequest := &types.ScheduleRequest{}
	err := jsonutils.UnmarshalBody(request.Body, scheduleRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	scheduleResponse, err := c.createScheduleUseCase.Execute(ctx, scheduleRequest, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(scheduleResponse)
}

// @Summary Creates a new schedule
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 404
// @Failure 500
// @Router /schedules/{uuid} [get]
func (c *SchedulesController) Get(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	uuid := mux.Vars(request)["uuid"]

	scheduleResponse, err := c.getScheduleUseCase.Execute(ctx, uuid, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(scheduleResponse)
}
