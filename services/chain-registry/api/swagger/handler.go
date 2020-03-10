package swagger

import (
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/rest"

	"github.com/gorilla/mux"
)

type Handler struct {
	serveSpecs string
	serveUI    string
}

func NewHandler(serveSpecs, serveUI string) *Handler {
	return &Handler{
		serveSpecs: serveSpecs,
		serveUI:    serveUI,
	}
}

func (h *Handler) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/swagger/swagger.json").Handler(rest.ServeFile(h.serveSpecs))
	router.Methods(http.MethodGet).PathPrefix("/swagger").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir(h.serveUI))))
}
