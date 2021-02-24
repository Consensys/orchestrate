package loadbalancer

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ConsenSys/orchestrate/pkg/http/config/dynamic"
	"github.com/sirupsen/logrus"
	"github.com/vulcand/oxy/roundrobin"
)

const cookieNameLength = 6

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error) {
	cfg, ok := configuration.(*dynamic.LoadBalancer)
	if !ok {
		return nil, nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	if len(cfg.Servers) == 0 {
		return nil, nil, fmt.Errorf("no server provided")
	}

	var urls []*url.URL
	for _, srv := range cfg.Servers {
		if srv.Weight < 0 {
			return nil, nil, fmt.Errorf("server weight should be >= 0 but got %v", srv.Weight)
		}

		var u *url.URL
		u, err = url.Parse(srv.URL)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing server URL %s: %v", srv.URL, err)
		}

		urls = append(urls, u)
	}

	var options []roundrobin.LBOption

	// Append logger option (it is used to deactivate heavy roundrobin logs at Debug Level)
	lbLogger := logrus.New()
	lbLogger.SetLevel(logrus.InfoLevel)
	lbLogger.SetFormatter(logrus.StandardLogger().Formatter)
	options = append(options, roundrobin.RoundRobinLogger(lbLogger))

	// Append Sticky option
	if cfg.Sticky != nil && cfg.Sticky.Cookie != nil {
		var cookieName string
		if cfg.Sticky.Cookie.Name == "" {
			cookieName, err = b.GenerateCookieName(name)
			if err != nil {
				return nil, nil, err
			}
		} else {
			cookieName = SanitizeCookieName(cfg.Sticky.Cookie.Name)
		}

		opts := roundrobin.CookieOptions{
			HTTPOnly: cfg.Sticky.Cookie.HTTPOnly,
			Secure:   cfg.Sticky.Cookie.Secure,
			Path:     cfg.Sticky.Cookie.Path,
			Domain:   cfg.Sticky.Cookie.Domain,
			MaxAge:   cfg.Sticky.Cookie.MaxAge,
		}

		options = append(options, roundrobin.EnableStickySession(roundrobin.NewStickySessionWithOptions(cookieName, opts)))
	}

	return func(h http.Handler) http.Handler {
		lb, _ := roundrobin.New(h, options...)
		for i, server := range cfg.Servers {
			_ = lb.UpsertServer(urls[i], roundrobin.Weight(server.Weight))
		}
		return lb
	}, nil, nil
}

func (b *Builder) GenerateCookieName(name string) (string, error) {
	data := []byte("_ORCHESTRATE_" + name)

	hash := sha1.New()
	_, err := hash.Write(data)
	if err != nil {
		// Impossible case
		return "", err
	}

	return fmt.Sprintf("_%x", hash.Sum(nil))[:cookieNameLength], nil
}

// SanitizeName According to [RFC 2616](https://www.ietf.org/rfc/rfc2616.txt) section 2.2
func SanitizeCookieName(name string) string {
	return strings.Map(cookieNameSanitizer, name)
}

func cookieNameSanitizer(r rune) rune {
	switch r {
	case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '`', '|', '~':
		return r
	}

	switch {
	case 'a' <= r && r <= 'z': //nolint
		fallthrough //nolint
	case 'A' <= r && r <= 'Z': //nolint
		fallthrough //nolint
	case '0' <= r && r <= '9': //nolint
		return r
	default:
		return '_'
	}
}
