package secure

import (
	"context"
	"fmt"
	"net/http"

	"github.com/unrolled/secure"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(ctx context.Context, _ string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	cfg, ok := configuration.(*dynamic.SecureHeaders)
	if !ok {
		return nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	m := secure.New(secure.Options{
		BrowserXssFilter:        cfg.BrowserXSSFilter,
		ContentTypeNosniff:      cfg.ContentTypeNosniff,
		ForceSTSHeader:          cfg.ForceSTSHeader,
		FrameDeny:               cfg.FrameDeny,
		IsDevelopment:           cfg.IsDevelopment,
		SSLRedirect:             cfg.SSLRedirect,
		SSLForceHost:            cfg.SSLForceHost,
		SSLTemporaryRedirect:    cfg.SSLTemporaryRedirect,
		STSIncludeSubdomains:    cfg.STSIncludeSubdomains,
		STSPreload:              cfg.STSPreload,
		ContentSecurityPolicy:   cfg.ContentSecurityPolicy,
		CustomBrowserXssValue:   cfg.CustomBrowserXSSValue,
		CustomFrameOptionsValue: cfg.CustomFrameOptionsValue,
		PublicKey:               cfg.PublicKey,
		ReferrerPolicy:          cfg.ReferrerPolicy,
		SSLHost:                 cfg.SSLHost,
		AllowedHosts:            cfg.AllowedHosts,
		HostsProxyHeaders:       cfg.HostsProxyHeaders,
		SSLProxyHeaders:         cfg.SSLProxyHeaders,
		STSSeconds:              cfg.STSSeconds,
		FeaturePolicy:           cfg.FeaturePolicy,
	})

	if cfg.IsProxy {
		return m.HandlerForRequestOnly, m.ModifyResponseHeaders, nil
	}

	return m.Handler, nil, nil
}
