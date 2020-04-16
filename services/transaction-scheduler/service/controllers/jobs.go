package controllers

import (
	"net/http"

	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"

	"github.com/gorilla/mux"
)

type JobsController struct {
	usecases *usecases.UseCases
}

func NewJobsController(uc *usecases.UseCases) *JobsController {
	return &JobsController{
		usecases: uc,
	}
}

// Add routes to router
func (c *JobsController) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/jobs").HandlerFunc(c.GetJobs)
}

// @Summary Retrieves a list of all transaction jobs matching a filter
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 400
// @Failure 500
// @Router /jobs [get]
func (c *JobsController) GetJobs(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	// TODO: Implement logic
}
