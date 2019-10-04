package grpcserver

import "google.golang.org/grpc"

// Enhancer are functions that enhance net/http Multiplexers
type Enhancer func(*grpc.Server) *grpc.Server

// ApplyEnhancers apply enhancers on a server
func ApplyEnhancers(s *grpc.Server, enhancers ...Enhancer) {
	// Enhance server
	for _, enhancer := range enhancers {
		enhancer(s)
	}
}
