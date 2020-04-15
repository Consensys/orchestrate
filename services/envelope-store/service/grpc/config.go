package grpc

import "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"

func NewStaticConfig() *static.Configuration {
	return &static.Configuration{
		Services: &static.Services{
			Envelopes: &static.Envelopes{},
		},
		Interceptors: []*static.Interceptor{
			{Tags: &static.Tags{}},
			{Logrus: &static.Logrus{}},
			{Auth: &static.Auth{}},
			{Error: &static.Error{}},
			{Recovery: &static.Recovery{}},
		},
	}
}
