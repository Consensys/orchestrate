package app

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/mock"
	mockprovider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/provider/mock"
	grpcstatic "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/config/static"
	mockinterceptor "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/interceptor/mock"
	mockservice "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/service/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
	mockhandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler/mock"
	mockmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/middleware/mock"
	dynrouter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/router/dynamic"
)

func TestProviderOpt(t *testing.T) {
	prvdr := mockprovider.New()

	app, err := New(newTestConfig())
	assert.NoError(t, err, "Creating app should not error")

	err = ProviderOpt(prvdr)(app)
	assert.NoError(t, err, "Applying option should not error")

	var listened bool
	app.AddListener(func(context.Context, interface{}) error {
		listened = true
		return nil
	})

	err = app.Start(context.Background())
	require.NoError(t, err, "App should have started properly")

	// Wait for application to properly start
	time.Sleep(100 * time.Millisecond)

	// Provide a message
	_ = prvdr.ProvideMsg(context.Background(), dynamic.NewMessage("test", &dynamic.Configuration{}))

	// Wait for application to properly process provided message
	time.Sleep(100 * time.Millisecond)

	// Stop app
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = app.Stop(ctx)
	assert.NoError(t, err, "App should have stop properly")

	assert.True(t, listened, "Listener should have been called")
}

func TestMiddlewareOpt(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	midBuilder := mockmid.NewMockBuilder(ctrlr)

	app, err := New(newTestConfig())
	assert.NoError(t, err, "Creating app should not error")

	err = MiddlewareOpt(
		reflect.TypeOf(&dynamic.Mock{}),
		midBuilder,
	)(app)
	assert.NoError(t, err, "Applying option should not error")

	testCfg := &dynamic.Mock{}
	cfg := &dynamic.Middleware{
		Mock: testCfg,
	}
	midBuilder.EXPECT().Build(gomock.Any(), "test", testCfg)
	_, _, err = app.HTTP().(*dynrouter.Builder).Middleware.Build(context.Background(), "test", cfg)
	assert.NoError(t, err, "Building middleware should not error")
}

func TestHandlerOpt(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	handlerBuilder := mockhandler.NewMockBuilder(ctrlr)

	app, err := New(newTestConfig())
	assert.NoError(t, err, "Creating app should not error")

	err = HandlerOpt(
		reflect.TypeOf(&dynamic.Mock{}),
		handlerBuilder,
	)(app)
	assert.NoError(t, err, "Applying option should not error")

	testCfg := &dynamic.Mock{}
	cfg := &dynamic.Service{
		Mock: testCfg,
	}
	handlerBuilder.EXPECT().Build(gomock.Any(), "test", testCfg, nil)
	_, err = app.HTTP().(*dynrouter.Builder).Handler.Build(context.Background(), "test", cfg, nil)
	assert.NoError(t, err, "Building middleware should not error")
}

func TestInterceptorOpt(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	interceptorBuilder := mockinterceptor.NewMockBuilder(ctrlr)
	opt := InterceptorOpt(
		reflect.TypeOf(&grpcstatic.Mock{}),
		interceptorBuilder,
	)

	testCfg := newTestConfig()
	mockCfg := &grpcstatic.Mock{}
	testCfg.GRPC.Static.Interceptors = []*grpcstatic.Interceptor{
		&grpcstatic.Interceptor{Mock: mockCfg},
	}

	app, err := New(testCfg, opt)
	assert.NoError(t, err, "Creating app should not error")

	interceptorBuilder.EXPECT().Build(gomock.Any(), gomock.Any(), mockCfg)
	_, err = app.GRPC().Build(context.Background(), "test", testCfg.GRPC.Static)
	assert.NoError(t, err, "Building interceptor should not error")
}

func TestServiceOpt(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	serviceBuilder := mockservice.NewMockBuilder(ctrlr)
	opt := ServiceOpt(
		reflect.TypeOf(&grpcstatic.Mock{}),
		serviceBuilder,
	)

	testCfg := newTestConfig()
	mockCfg := &grpcstatic.Mock{}
	testCfg.GRPC.Static.Services = &grpcstatic.Services{Mock: mockCfg}

	app, err := New(testCfg, opt)
	assert.NoError(t, err, "Creating app should not error")

	serviceBuilder.EXPECT().Build(gomock.Any(), gomock.Any(), mockCfg)
	_, err = app.GRPC().Build(context.Background(), "test", testCfg.GRPC.Static)
	assert.NoError(t, err, "Building interceptor should not error")
}

func TestMultitenancyOpt(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	jwtChecker := mockauth.NewMockChecker(ctrlr)
	keyChecker := mockauth.NewMockChecker(ctrlr)

	opt := MultiTenancyOpt("auth", jwtChecker, keyChecker, true)

	_, err := New(newTestConfig(), opt)
	assert.NoError(t, err, "Creating app should not error")
}

func TestLoggerOpt(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	opt := LoggerOpt("test")

	_, err := New(newTestConfig(), opt)
	assert.NoError(t, err, "Creating app should not error")
}

func TestSwaggerOpt(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	opt := SwaggerOpt("test")

	_, err := New(newTestConfig(), opt)
	assert.NoError(t, err, "Creating app should not error")
}

func TestMetricsOpt(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	opt := MetricsOpt()

	testCfg := newTestConfig()
	_, err := New(testCfg, opt)
	assert.NoError(t, err, "Creating app should not error")
}
