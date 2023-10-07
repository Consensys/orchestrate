package headers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/middleware"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/middleware/secure"
	"github.com/justinas/alice"
)

type Builder struct {
	secure middleware.Builder
}

func NewBuilder() *Builder {
	return &Builder{
		secure: secure.NewBuilder(),
	}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	cfg, ok := configuration.(*dynamic.Headers)
	if !ok {
		return nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	chain := alice.New()
	var respModifiers []func(*http.Response) error
	if cfg.Secure != nil {
		cfg.Secure.IsProxy = cfg.IsProxy
		mid, respModifier, err = b.secure.Build(ctx, name, configuration)
		if err != nil {
			return nil, nil, err
		}
		chain.Append(mid)
		respModifiers = append(respModifiers, respModifier)
	}

	m := New(cfg.Custom)

	if cfg.IsProxy {
		return m.HandlerForRequestOnly, httputil.CombineResponseModifiers(append(respModifiers, m.ModifyResponseHeaders)...), nil
	}

	return m.Handler, httputil.CombineResponseModifiers(respModifiers...), nil
}

type Headers struct {
	cfg *dynamic.CustomHeaders
}

func New(cfg *dynamic.CustomHeaders) *Headers {
	return &Headers{
		cfg: cfg,
	}
}

// Handler implements the http.HandlerFunc for integration with the standard net/http lib.
func (s *Headers) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		s.modifyCustomRequestHeaders(req)
		s.modifyResponseHeaders(rw)
		h.ServeHTTP(rw, req)
	})
}

// HandlerForRequestOnly implements the http.HandlerFunc for integration with the standard net/http lib.
// Note that this is for requests only and will not write any headers.
func (s *Headers) HandlerForRequestOnly(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		s.modifyCustomRequestHeaders(req)
		h.ServeHTTP(rw, req)
	})
}

func (s *Headers) modifyCustomRequestHeaders(req *http.Request) {
	// Loop through Custom request headers
	for header, value := range s.cfg.RequestHeaders {
		switch {
		case value == "":
			req.Header.Del(header)

		case strings.EqualFold(header, "Host"):
			req.Host = value

		default:
			req.Header.Set(header, value)
		}
	}
}

// addResponseHeaders Adds the headers from 'responseHeader' to the response.
func (s *Headers) modifyResponseHeaders(w http.ResponseWriter) {
	for key, value := range s.cfg.ResponseHeaders {
		if value == "" {
			w.Header().Del(key)
		} else {
			w.Header().Set(key, value)
		}
	}
}

// ModifyResponseHeaders set or delete response headers.
// This method is called AFTER the response is generated from the backend
// and can merge/override headers from the backend response.
func (s *Headers) ModifyResponseHeaders(res *http.Response) error {
	// Loop through Custom response headers
	for header, value := range s.cfg.ResponseHeaders {
		if value == "" {
			res.Header.Del(header)
		} else {
			res.Header.Set(header, value)
		}
	}

	return nil
}
