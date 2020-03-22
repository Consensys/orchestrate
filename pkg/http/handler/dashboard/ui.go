package dashboard

import (
	"net/http"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/gorilla/mux"
)

// UI exposes dashboard user interface
type UI struct {
	files http.Handler
}

func NewUI(files http.Handler) *UI {
	return &UI{
		files: files,
	}
}

// Append add dashboard routes on a router
func (h *UI) Append(router *mux.Router) {
	if h.files == nil {
		log.WithoutContext().Error("No dashboard UI to serve")
		return
	}

	// Expose dashboard
	router.Methods(http.MethodGet).
		Path("/").
		HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			http.Redirect(rw, req, req.Header.Get("X-Forwarded-Prefix")+"/dashboard/", http.StatusFound)
		})

	router.Methods(http.MethodGet).
		PathPrefix("/dashboard/").
		Handler(http.StripPrefix("/dashboard/", h.files))
}

func (h *UI) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.files.ServeHTTP(rw, req)
}
