package prometheus

var testCfg = &Config{
	TCP:  &TCPConfig{Buckets: []float64{0.1, 0.3, 1.2, 5.0}},
	HTTP: &HTTPConfig{Buckets: []float64{0.1, 0.3, 1.2, 5.0}},
	GRPC: &GRPCServerConfig{Buckets: []float64{0.1, 0.3, 1.2, 5.0}},
}
