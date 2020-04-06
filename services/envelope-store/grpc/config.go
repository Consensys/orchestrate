package grpc

import "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"

func NewStaticConfig() *static.Configuration {
	return &static.Configuration{
		Services: &static.Services{
			Envelopes: &static.Envelopes{},
		},
		Interceptors: []*static.Interceptor{
			&static.Interceptor{Tags: &static.Tags{}},
			&static.Interceptor{Logrus: &static.Logrus{}},
			&static.Interceptor{Auth: &static.Auth{}},
			&static.Interceptor{Error: &static.Error{}},
			&static.Interceptor{Recovery: &static.Recovery{}},
		},
	}
}