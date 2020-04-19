package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Prometheus struct {
	tcp        *TCP
	http       *HTTP
	grpcServer *GRPCServer
}

func New(cfg *Config) *Prometheus {
	return &Prometheus{
		tcp:        NewTCP(cfg),
		http:       NewHTTP(cfg),
		grpcServer: NewGRPCServer(cfg),
	}
}

func (reg *Prometheus) TCP() *TCP {
	return reg.tcp
}

func (reg *Prometheus) HTTP() *HTTP {
	return reg.http
}

func (reg *Prometheus) GRPCServer() *GRPCServer {
	return reg.grpcServer
}

// Describe implements prometheus.Collector and simply calls
// the registered describer functions.
func (reg *Prometheus) Describe(ch chan<- *prometheus.Desc) {
	reg.tcp.Describe(ch)
	reg.http.Describe(ch)
	reg.grpcServer.Describe(ch)
}

// Collect collectors
func (reg *Prometheus) Collect(ch chan<- prometheus.Metric) {
	reg.tcp.Collect(ch)
	reg.http.Collect(ch)
	reg.grpcServer.Collect(ch)
}
