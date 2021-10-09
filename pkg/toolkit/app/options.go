package app

import (
	"fmt"
	"math"
	"reflect"

	"github.com/consensys/orchestrate/pkg/toolkit/app/auth"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/provider"
	staticprovider "github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/provider/static"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/handler"
	dynhandler "github.com/consensys/orchestrate/pkg/toolkit/app/http/handler/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/handler/healthcheck"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/handler/prometheus"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/handler/swagger"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/middleware"
	authmid "github.com/consensys/orchestrate/pkg/toolkit/app/http/middleware/auth"
	dynmid "github.com/consensys/orchestrate/pkg/toolkit/app/http/middleware/dynamic"
	multitenancymid "github.com/consensys/orchestrate/pkg/toolkit/app/http/middleware/multitenancy"
	dynrouter "github.com/consensys/orchestrate/pkg/toolkit/app/http/router/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/metrics"
	metricregistry "github.com/consensys/orchestrate/pkg/toolkit/app/metrics/registry"
	healthz "github.com/heptiolabs/healthcheck"
	prom "github.com/prometheus/client_golang/prometheus"
	traefikdynamic "github.com/traefik/traefik/v2/pkg/config/dynamic"
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

	// Provider for Auth middleware dynamic configuration
	cfg := dynamic.NewConfig()
	cfg.HTTP.Middlewares[authMidName] = &dynamic.Middleware{
		Auth: &dynamic.Auth{},
	}
	providerOpt := ProviderOpt(staticprovider.New(dynamic.NewMessage("multitenancy", cfg)))

	return CombineOptions(
		authMidOpt,
		multitenancyMidOpt,
		providerOpt,
	)
}

func LoggerOpt(midName string) Option {
	return CombineOptions(
		LoggerMiddlewareOpt(midName),
	)
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
			EntryPoints: []string{http.DefaultHTTPAppEntryPoint},
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

func MetricsOpt(appMetrics ...metrics.Prometheus) Option {
	registry := prom.NewRegistry()

	// Register Provided metrics
	appMetricsRegister := func(app *App) error {
		// Register base Process and Golang runtime metrics
		if app.cfg.Metrics.IsActive(metricregistry.GoMetricsModule) {
			app.metricReg.Add(prom.NewGoCollector())
		}
		if app.cfg.Metrics.IsActive(metricregistry.ProcessMetricsModule) {
			app.metricReg.Add(prom.NewProcessCollector(prom.ProcessCollectorOpts{}))
		}

		for _, m := range appMetrics {
			app.metricReg.Add(m)
		}

		return nil
	}

	healthzOpt := func(app *App) error {
		var h healthz.Handler
		if app.cfg.Metrics.IsActive(metricregistry.HealthzMetricsModule) {
			h = healthz.NewMetricsHandler(registry, metrics.Namespace)
		} else {
			h = healthz.NewHandler()
		}

		return HealthcheckOpt(h)(app)
	}

	return CombineOptions(
		appMetricsRegister,
		healthzOpt,
		PrometheusOpt(registry),
		DashboardOpt(),
	)
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

func HealthcheckOpt(h healthz.Handler, middlewares ...string) Option {
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

	return func(app *App) error {
		for _, chk := range app.readinessCheck {
			h.AddReadinessCheck(chk.Name, chk.Check)
		}

		// Handler builder option
		handlerOpt := HandlerOpt(
			reflect.TypeOf(&dynamic.HealthCheck{}),
			healthcheck.NewBuilder(h),
		)

		return CombineOptions(
			providerOpt,
			handlerOpt,
		)(app)
	}
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

	// Register Prometheus registry
	promRegister := func(app *App) error {
		// Register app metric collector
		err := registry.Register(app.MetricRegistry())
		if err != nil {
			return err
		}

		return nil
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
