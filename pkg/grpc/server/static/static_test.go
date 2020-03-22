// +build unit

package static

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	interceptormock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/mock"
	servermock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server/mock"
	servicemock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/service/mock"
)

func TestBuilder(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	interceptorMock := interceptormock.NewMockBuilder(ctrlr)
	serviceMock := servicemock.NewMockBuilder(ctrlr)
	optionsMock := servermock.NewMockOptionsBuilder(ctrlr)
	b := NewBuilder()

	b.Options = optionsMock
	b.Interceptor = interceptorMock
	b.Service = serviceMock

	name := "test"
	cfg := &static.Configuration{
		Options: &static.Options{},
		Interceptors: []*static.Interceptor{
			&static.Interceptor{},
			&static.Interceptor{},
		},
		Services: &static.Services{},
	}

	call1 := interceptorMock.EXPECT().Build(gomock.Any(), name, cfg.Interceptors[0]).Return(nil, nil, nil, nil)
	interceptorMock.EXPECT().Build(gomock.Any(), name, cfg.Interceptors[1]).Return(nil, nil, nil, nil).After(call1)
	serviceMock.EXPECT().Build(gomock.Any(), name, cfg.Services).Return(nil, nil)
	optionsMock.EXPECT().Build(gomock.Any(), name, cfg.Options).Return(nil, nil)

	srv, err := b.Build(context.Background(), name, cfg)
	require.NoError(t, err)
	assert.NotNil(t, srv, "Server should have been created")
}
