package dynamic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/runtime"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/provider"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/handler"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/handler/dashboard"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/middleware"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/middleware/accesslog"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/router"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	tlsmanager "github.com/consensys/orchestrate/pkg/toolkit/tls/manager"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
	traefikstatic "github.com/traefik/traefik/v2/pkg/config/static"
	tlog "github.com/traefik/traefik/v2/pkg/log"
	"github.com/traefik/traefik/v2/pkg/middlewares/requestdecorator"
	"github.com/traefik/traefik/v2/pkg/rules"
)

type Builder struct {
	Middleware middleware.Builder
	Handler    handler.Builder
	TLS        tlsmanager.Manager
	Metrics    middleware.Builder

	dashboard handler.Builder

	accesslog    middleware.Builder
	epaccesslogs map[string]func(http http.Handler) http.Handler

	reqdecorator *requestdecorator.RequestDecorator
	logger       *log.Logger
}

func NewBuilder(staticCfg *traefikstatic.Configuration) *Builder {
	b := &Builder{
		dashboard:    dashboard.NewBuilder(staticCfg),
		accesslog:    accesslog.NewBuilder(),
		epaccesslogs: make(map[string]func(http http.Handler) http.Handler),
		reqdecorator: requestdecorator.New(staticCfg.HostResolver),
		logger:       log.NewLogger().SetComponent("router"),
	}

	return b
}

func (b *Builder) Build(ctx context.Context, entryPointNames []string, configuration interface{}) (map[string]*router.Router, error) {
	cfg, ok := configuration.(*dynamic.Configuration)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	infos := runtime.NewInfos(cfg)
	infos.PopulateUsedBy()

	routers := make(map[string]*router.Router)
	for _, entryPointName := range entryPointNames {
		routers[entryPointName] = &router.Router{}
	}

	err := b.buildRouters(ctx, routers, infos, false)
	if err != nil {
		return nil, err
	}

	err = b.buildRouters(ctx, routers, infos, true)
	if err != nil {
		return nil, err
	}

	return routers, nil
}

func (b *Builder) buildRouters(ctx context.Context, routers map[string]*router.Router, infos *runtime.Infos, isTLS bool) error {
	var entryPointNames []string
	for entryPointName := range routers {
		entryPointNames = append(entryPointNames, entryPointName)
	}

	for entryPointName, rtInfos := range infos.RouterInfosByEntryPoint(ctx, entryPointNames, isTLS) {
		logger := b.logger.WithField("entrypoint", entryPointName)
		epCtx := tlog.With(
			httputil.WithEntryPoint(ctx, entryPointName),
			tlog.Str("entrypoint", entryPointName),
		)
		logger.Debug("building entrypoint router")

		mux, err := rules.NewRouter()
		if err != nil {
			return err
		}

		epAccessLogMiddleware, ok := b.epaccesslogs[entryPointName]
		if ok {
			logger.Debugf("accesslog activated on entrypoint")
		}

		for routerName, rtInfo := range rtInfos {
			logger := b.logger.WithField("router_name", routerName)
			rtCtx := tlog.With(
				httputil.WithRouter(provider.WithName(epCtx, routerName), routerName),
				tlog.Str("router_name", routerName),
			)

			logger.WithFields(logrus.Fields{
				"rule.name":    rtInfo.Router.Rule,
				"middlewares":  rtInfo.Router.Middlewares,
				"service.name": rtInfo.Router.Service,
				"priority":     rtInfo.Router.Priority,
			}).Debug("building route")

			var h http.Handler
			h, err = b.buildHandler(rtCtx, routerName, rtInfo, infos, epAccessLogMiddleware)
			if err != nil {
				rtInfo.AddError(err, true)
				logger.WithError(err).Error("could not build route handler")
				continue
			}

			err = mux.AddRoute(rtInfo.Rule, rtInfo.Priority, h)
			if err != nil {
				rtInfo.AddError(err, true)
				logger.WithError(err).Error("could not add route on router")
				continue
			}
		}

		conf, ok := routers[entryPointName]
		if !ok {
			continue
		}

		// RequestDecorator middleware is necessary for using traefik routing
		finalHandler, err := requestdecorator.WrapHandler(b.reqdecorator)(mux)
		if err != nil {
			return err
		}

		if isTLS {
			conf.HTTPS = finalHandler
			if b.TLS != nil {
				conf.TLSConfig, conf.HostTLSConfigs, err = b.TLS.Get(ctx, rtInfos)
				if err != nil {
					return err
				}
			}
		} else {
			conf.HTTP = finalHandler
		}
	}

	return nil
}

func (b *Builder) buildHandler(ctx context.Context, routerName string, rtInfo *runtime.RouterInfo, infos *runtime.Infos, accessLog func(http.Handler) http.Handler) (http.Handler, error) {
	hCtx := httputil.WithService(ctx, rtInfo.Service)

	mid, respModifier, rvErr := b.buildMiddleware(
		hCtx,
		routerName,
		rtInfo,
		infos,
		accessLog,
	)

	h, err := b.buildService(
		hCtx,
		fmt.Sprintf("%v:%v", routerName, rtInfo.Service),
		infos.Services[rtInfo.Service],
		infos,
		respModifier,
	)
	if err != nil {
		rvErr = err
	}

	return mid(h), rvErr
}

func (b *Builder) buildMiddleware(ctx context.Context, routerName string, rtInfo *runtime.RouterInfo, infos *runtime.Infos, accessLog func(http.Handler) http.Handler) (func(http.Handler) http.Handler, func(*http.Response) error, error) { // nolint
	chain := alice.New()
	var respModifiers []func(resp *http.Response) error
	var rvErr error

	// Add metrics middleware
	if b.Metrics != nil {
		mid, respModifier, err := b.Metrics.Build(
			ctx,
			fmt.Sprintf("%v:%v", routerName, "metrics"),
			nil,
		)
		if err != nil {
			b.logger.WithError(err).Error("could not build metrics middleware")
			rvErr = err
		} else {
			if mid != nil {
				chain = chain.Append(mid)
			}

			if respModifier != nil {
				respModifiers = append(respModifiers, respModifier)
			}
		}
	}

	for _, midName := range rtInfo.Middlewares {
		midCtx := httputil.WithMiddleware(ctx, midName)
		logger := b.logger.WithContext(midCtx).WithField("middleware", midName)

		// In case a services is missing one of the middleware configurationg we skip it usage and warning
		if infos.Middlewares[midName] == nil {
			continue
		}

		if infos.Middlewares[midName].Middleware == nil {
			rvErr = fmt.Errorf("middleware %q configuration is empty", midName)
			logger.Error("failed to build middleware with empty configuration")
			rtInfo.AddError(rvErr, true)
			continue
		}

		switch {
		case infos.Middlewares[midName].Middleware.AccessLog != nil:
			// Treat particular case of access logs
			mid, _, err := b.accesslog.Build(
				midCtx,
				midName,
				infos.Middlewares[midName].Middleware.AccessLog,
			)
			if err != nil {
				infos.Middlewares[midName].AddError(err, true)
				logger.WithError(err).Error("could not build middleware")
				rvErr = err
				continue
			}

			// Add accesslog middleware to the chain
			chain = chain.Append(mid)

			// Set accessLog to nil to make sure we do not register
			// accesslog middleware twice
			accessLog = nil
			continue
		case b.Middleware != nil:
			mid, respModifier, err := b.Middleware.Build(
				midCtx,
				fmt.Sprintf("%v:%v", routerName, midName),
				infos.Middlewares[midName].Middleware,
			)
			if err != nil {
				infos.Middlewares[midName].AddError(err, true)
				logger.WithError(err).Error("could not build middleware")
				rvErr = err
				continue
			}

			if mid != nil {
				chain = chain.Append(mid)
			}

			if respModifier != nil {
				respModifiers = append(respModifiers, respModifier)
			}
		default:
			logger.Debugf("no middleware builder registered")
		}
	}

	if accessLog != nil {
		b.logger.Debugf("added entrypoint accesslog")
		chain = alice.New(accessLog).Extend(chain)
	}

	return chain.Then, httputil.CombineResponseModifiers(respModifiers...), rvErr
}

func (b *Builder) buildService(ctx context.Context, serviceName string, srvInfo *runtime.ServiceInfo, infos *runtime.Infos, respModifier func(*http.Response) error) (http.Handler, error) {
	logger := b.logger.WithField("service_name", serviceName)

	switch {
	case srvInfo.Service.Dashboard != nil:
		h, err := b.dashboard.Build(
			ctx,
			serviceName,
			infos,
			nil,
		)
		if err != nil {
			srvInfo.AddError(err, true)
			logger.WithError(err).Error("could not build handler")
			return http.NotFoundHandler(), err
		}
		return h, nil
	case b.Handler != nil:
		h, err := b.Handler.Build(
			ctx,
			serviceName,
			srvInfo.Service,
			respModifier,
		)
		if err != nil {
			srvInfo.AddError(err, true)
			logger.WithError(err).Error("could not build handler")
			return http.NotFoundHandler(), err
		}
		return h, nil
	default:
		logger.Debug("no handler builder registered")
		return http.NotFoundHandler(), fmt.Errorf("no handler to build (falling back on NotFound)")
	}
}
