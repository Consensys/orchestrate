package dashboard

import (
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime"

	"github.com/gorilla/mux"
)

func init() {
	// FIXME Goroutines2 -> Goroutines
	expvar.Publish("Goroutines", expvar.Func(goroutines))
}

func goroutines() interface{} {
	return runtime.NumGoroutine()
}

// DebugHandler expose debug routes
type Debug struct{}

func NewDebug() *Debug {
	return &Debug{}
}

// Append add debug routes on a router
func (h *Debug) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/debug/vars").HandlerFunc(h.ServeHTTPVars)

	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(5)
	router.Methods(http.MethodGet).PathPrefix("/debug/pprof/cmdline").HandlerFunc(pprof.Cmdline)
	router.Methods(http.MethodGet).PathPrefix("/debug/pprof/profile").HandlerFunc(pprof.Profile)
	router.Methods(http.MethodGet).PathPrefix("/debug/pprof/symbol").HandlerFunc(pprof.Symbol)
	router.Methods(http.MethodGet).PathPrefix("/debug/pprof/trace").HandlerFunc(pprof.Trace)
	router.Methods(http.MethodGet).PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)
}

func (h *Debug) ServeHTTPVars(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(rw, "{\n")
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprint(rw, ",\n")
		}
		first = false
		fmt.Fprintf(rw, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprint(rw, "\n}\n")
}
