package controllers

import (
	"github.com/gorilla/mux"
)

type EthereumController struct{}

func NewEthereumController() *EthereumController {
	return &EthereumController{}
}

// Add routes to router
func (c *EthereumController) Append(router *mux.Router) {}
