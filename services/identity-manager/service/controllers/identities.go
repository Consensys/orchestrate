package controllers

import (
	"github.com/gorilla/mux"
)

type IdentitiesController struct{}

func NewIdentitiesController() *IdentitiesController {
	return &IdentitiesController{}
}

// Add routes to router
func (c *IdentitiesController) Append(router *mux.Router) {}
