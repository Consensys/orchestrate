package app

import (
	"fmt"
	"math"
	"reflect"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	staticprovider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor"
	grpcauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/auth"
	grpclogrus "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/logrus"
	grpcmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/metrics"
	staticinterceptor "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/static"
	staticgrpc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/service"
	staticservice "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/service/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler"
	dynhandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/prometheus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/swagger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware"
	authmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/auth"
	dynmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/dynamic"
	multitenancymid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/multitenancy"
	dynrouter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router/dynamic"
	multimetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/multi"
)

type Option func(*App) error

func ProviderOpt(prvdr provider.Provider) Option {
	return func(app *App) error {
		app.AddProvider(prvdr)
		return nil
	}
}

func MiddlewareOpt(typ reflect.Type, builder middleware.Builder) Option {
	return func(app *App) error {
		httpBuilder, ok := app.HTTP().(*dynrouter.Builder)
		if ok {
			midBuilder, ok := httpBuilder.Middleware.(*dynmid.Builder)
			if ok {
				midBuilder.AddBuilder(typ, builder)
				return nil
			}
			return fmt.Errorf("invalid router Middleware builder of type %T (expected %T)", httpBuilder.Middleware, midBuilder)
		}
		return fmt.Errorf("can not add Middleware builder on router of type %T (expected %T)", app.HTTP(), httpBuilder)
	}
}

func HandlerOpt(typ reflect.Type, builder handler.Builder) Option {
	return func(app *App) error {
		httpBuilder, ok := app.HTTP().(*dynrouter.Builder)
		if ok {
			handlerBuilder, ok := httpBuilder.Handler.(*dynhandler.Builder)
			if ok {
				handlerBuilder.AddBuilder(typ, builder)
				return nil
			}
			return fmt.Errorf("invalid router Handler builder of type %T (expected %T)", httpBuilder.Handler, handlerBuilder)

		}
		return fmt.Errorf("can not add http.Handler builder on router of type %T (expect %T)", app.HTTP(), httpBuilder)
	}
}

func InterceptorOpt(typ reflect.Type, builder interceptor.Builder) Option {
	return func(app *App) error {
		grpcBuilder, ok := app.GRPC().(*staticgrpc.Builder)
		if ok {
			grpcBuilder.Interceptor.(*staticinterceptor.Builder).AddBuilder(typ, builder)
		}
		return nil
	}
}

func ServiceOpt(typ reflect.Type, builder service.Builder) Option {
	return func(app *App) error {
		grpcBuilder, ok := app.GRPC().(*staticgrpc.Builder)
		if ok {
			grpcBuilder.Service.(*staticservice.Builder).AddBuilder(typ, builder)
		}
		return nil
	}
}

func MultiTenancyOpt(authMidName string, jwt, key auth.Checker, multitenancy bool) Option {
	// Auth HTTP middleware
	authMidOpt := MiddlewareOpt(
		reflect.TypeOf(&dynamic.Auth{}),
		authmid.NewBuilder(jwt, key, multitenancy),
	)

	// Multitenancy HTTP middleware
	multitenancyMidOpt := MiddlewareOpt(
		reflect.TypeOf(&dynamic.MultiTenancy{}),
		multitenancymid.NewBuilder(),
	)

	// Auth GRPC Interceptor
	interceptorOpt := InterceptorOpt(
		reflect.TypeOf(&static.Auth{}),
		grpcauth.NewBuilder(auth.CombineCheckers(key, jwt), multitenancy),
	)

	// Provider for Auth middleware dynamic configuration
	cfg := dynamic.NewConfig()
	cfg.HTTP.Middlewares[authMidName] = &dynamic.Middleware{
		Auth: &dynamic.Auth{},
	}
	providerOpt := ProviderOpt(staticprovider.New(dynamic.NewMessage("multitenancy", cfg)))

	return CombineOptions(
		authMidOpt,
		multitenancyMidOpt,
		interceptorOpt,
		providerOpt,
	)
}

func LoggerOpt(midName string) Option {
	return CombineOptions(
		LoggerInterceptorOpt(),
		LoggerMiddlewareOpt(midName),
	)
}

func LoggerInterceptorOpt() Option {
	return func(app *App) error {
		return InterceptorOpt(
			reflect.TypeOf(&static.Logrus{}),
			grpclogrus.NewBuilder(app.Logger(), logrus.Fields{"system": "grpc.internal"}),
		)(app)
	}
}

func LoggerMiddlewareOpt(midName string) Option {
	// Provider injecting dynamic middleware configuration
	return func(app *App) error {
		cfg := dynamic.NewConfig()

		logFormat := ""
		if app.cfg.Log != nil {
			logFormat = app.cfg.Log.Format
		}

		cfg.HTTP.Middlewares[midName] = &dynamic.Middleware{
			AccessLog: &dynamic.AccessLog{
				Format: logFormat,
			},
		}

		return ProviderOpt(staticprovider.New(dynamic.NewMessage("logger-"+midName, cfg)))(app)
	}
}

func SwaggerOpt(specsFile string, middlewares ...string) Option {
	// Provider injecting dynamic middleware configuration
	cfg := dynamic.NewConfig()

	// Router to swagger
	cfg.HTTP.Routers["swagger"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultHTTPEntryPoint},
			Service:     "swagger",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/swagger`)",
			Middlewares: middlewares,
		},
	}

	// Swagger
	cfg.HTTP.Services["swagger"] = &dynamic.Service{
		Swagger: &dynamic.Swagger{
			SpecsFile: specsFile,
		},
	}

	providerOpt := ProviderOpt(staticprovider.New(dynamic.NewMessage("swagger", cfg)))

	// Option for Swagger handler
	handlerOpt := HandlerOpt(
		reflect.TypeOf(&dynamic.Swagger{}),
		swagger.NewBuilder(),
	)

	return CombineOptions(
		providerOpt,
		handlerOpt,
	)
}

func MetricsOpt(middlewares ...string) Option {
	registry := prom.NewRegistry()
	return CombineOptions(
		HealthcheckOpt(middlewares...),
		PrometheusOpt(registry, middlewares...),
		DashboardOpt(middlewares...),
		MetricsInterceptorOpt(),
	)
}

func MetricsInterceptorOpt() Option {
	return func(app *App) error {
		return InterceptorOpt(
			reflect.TypeOf(&static.Metrics{}),
			grpcmetrics.NewBuilder(app.Metrics().GRPCServer()),
		)(app)
	}
}

func DashboardOpt(middlewares ...string) Option {
	// Provider injecting dynamic middleware configuration
	cfg := dynamic.NewConfig()

	cfg.HTTP.Routers["dashboard"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultMetricsEntryPoint},
			Service:     "dashboard",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/api`) || PathPrefix(`/dashboard`)",
			Middlewares: append(middlewares, "strip-api"),
		},
	}

	cfg.HTTP.Middlewares["strip-api"] = &dynamic.Middleware{
		Middleware: &traefikdynamic.Middleware{
			StripPrefixRegex: &traefikdynamic.StripPrefixRegex{
				Regex: []string{"/api"},
			},
		},
	}

	// Dashboard API
	cfg.HTTP.Services["dashboard"] = &dynamic.Service{
		Dashboard: &dynamic.Dashboard{},
	}

	return ProviderOpt(staticprovider.New(dynamic.NewMessage("dashboard", cfg)))
}

func HealthcheckOpt(middlewares ...string) Option {
	// Provider injecting dynamic configuration
	cfg := dynamic.NewConfig()

	// Router to Healthchecks
	cfg.HTTP.Routers["healthcheck"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultMetricsEntryPoint},
			Service:     "healthcheck",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/live`) || PathPrefix(`/ready`)",
			Middlewares: middlewares,
		},
	}

	// Healthcheck
	cfg.HTTP.Services["healthcheck"] = &dynamic.Service{
		HealthCheck: &dynamic.HealthCheck{},
	}
	providerOpt := ProviderOpt(staticprovider.New(dynamic.NewMessage("healthcheck", cfg)))

	// Handler builder option
	handlerOpt := HandlerOpt(
		reflect.TypeOf(&dynamic.HealthCheck{}),
		healthcheck.NewTraefikBuilder(),
	)

	return CombineOptions(
		providerOpt,
		handlerOpt,
	)
}

func PrometheusOpt(registry *prom.Registry, middlewares ...string) Option {
	handlerOpt := HandlerOpt(
		reflect.TypeOf(&dynamic.Prometheus{}),
		prometheus.NewBuilder(registry),
	)

	// Provider injecting dynamic middleware configuration
	cfg := dynamic.NewConfig()
	cfg.HTTP.Routers["prometheus"] = &dynamic.Router{
		Router: &traefikdynamic.Router{
			EntryPoints: []string{http.DefaultMetricsEntryPoint},
			Service:     "prometheus",
			Priority:    math.MaxInt32,
			Rule:        "PathPrefix(`/metrics`)",
			Middlewares: middlewares,
		},
	}

	cfg.HTTP.Services["prometheus"] = &dynamic.Service{
		Prometheus: &dynamic.Prometheus{},
	}

	providerOpt := ProviderOpt(staticprovider.New(dynamic.NewMessage("prometheus", cfg)))

	// Register Prometheus regstry
	promRegister := func(app *App) error {
		// Register base Process and Golang runtime metrics
		registry.MustRegister(prom.NewProcessCollector(prom.ProcessCollectorOpts{}))
		registry.MustRegister(prom.NewGoCollector())

		// Register base custom metrics
		reg, ok := app.Metrics().(*multimetrics.Multi)
		if ok {
			return registry.Register(reg.Prometheus())
		}
		return fmt.Errorf("invalid metrics registry type %T (expected %T)", app.Metrics(), reg)
	}

	return CombineOptions(
		handlerOpt,
		providerOpt,
		promRegister,
	)
}

func CombineOptions(opts ...Option) Option {
	return func(app *App) error {
		for _, opt := range opts {
			err := opt(app)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
