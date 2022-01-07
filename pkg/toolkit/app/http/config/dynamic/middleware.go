package dynamic

import (
	"time"

	traefiktypes2 "github.com/traefik/paerser/types"
	traefikdynamic "github.com/traefik/traefik/v2/pkg/config/dynamic"
	traefiktypes "github.com/traefik/traefik/v2/pkg/types"

	"github.com/consensys/orchestrate/pkg/utils"
)

// +k8s:deepcopy-gen=true

type Middleware struct {
	*traefikdynamic.Middleware
	Auth         *Auth         `json:"auth,omitempty" toml:"auth,omitempty" yaml:"auth,omitempty"`
	MultiTenancy *MultiTenancy `json:"multitenancy,omitempty" toml:"multitenancy,omitempty" yaml:"multitenancy,omitempty"`
	Headers      *Headers      `json:"headers,omitempty" toml:"headers,omitempty" yaml:"headers,omitempty"`
	Cors         *Cors         `json:"cors,omitempty" toml:"cors,omitempty" yaml:"cors,omitempty"`
	LoadBalancer *LoadBalancer `json:"loadBalancer,omitempty" toml:"loadBalancer,omitempty" yaml:"loadBalancer,omitempty"`
	RateLimit    *RateLimit    `json:"rateLimit,omitempty" toml:"rateLimit,omitempty" yaml:"rateLimit,omitempty"`
	HTTPTrace    *HTTPTrace    `json:"httpTrace,omitempty" toml:"httpTrace,omitempty" yaml:"httpTrace,omitempty"`
	HTTPCache    *HTTPCache    `json:"httpCache,omitempty" toml:"httpCache,omitempty" yaml:"httpCache,omitempty"`
	Mock         *Mock         `json:"mock,omitempty" toml:"mock,omitempty" yaml:"mock,omitempty"`
	AccessLog    *AccessLog    `json:"accessLog,omitempty" toml:"accessLog,omitempty" yaml:"accessLog,omitempty"`
}

func (m *Middleware) Type() string {
	if m.Middleware != nil {
		return utils.ExtractType(m.Middleware)
	}

	return utils.ExtractType(m)
}

func (m *Middleware) Field() (interface{}, error) {
	return utils.ExtractField(m)
}

func FromTraefikMiddleware(middleware *traefikdynamic.Middleware) *Middleware {
	return &Middleware{
		Middleware: middleware,
	}
}

func ToTraefikMiddleware(middleware *Middleware) *traefikdynamic.Middleware {
	return middleware.Middleware
}

// +k8s:deepcopy-gen=true

type Auth struct{}

// +k8s:deepcopy-gen=true

type MultiTenancy struct {
	Tenant  string `json:"tenant,omitempty" toml:"tenant,omitempty" yaml:"tenant,omitempty"`
	OwnerID string `json:"owner_id,omitempty" toml:"tenant,omitempty" yaml:"tenant,omitempty"`
}

// +k8s:deepcopy-gen=true

type Headers struct {
	Secure *SecureHeaders `json:"secure,omitempty" toml:"secure,omitempty" yaml:"secure,omitempty"`
	Custom *CustomHeaders `json:"custom,omitempty" toml:"custom,omitempty" yaml:"custom,omitempty"`

	IsProxy bool `json:"isProxy,omitempty" toml:"isProxy,omitempty" yaml:"isProxy,omitempty"`
}

// +k8s:deepcopy-gen=true

type SecureHeaders struct {
	IsDevelopment           bool              `json:"isDevelopment,omitempty" toml:"isDevelopment,omitempty" yaml:"isDevelopment,omitempty"`
	IsProxy                 bool              `json:"isProxy,omitempty" toml:"isProxy,omitempty" yaml:"isProxy,omitempty"`
	SSLRedirect             bool              `json:"sslRedirect,omitempty" toml:"sslRedirect,omitempty" yaml:"sslRedirect,omitempty"`
	SSLTemporaryRedirect    bool              `json:"sslTemporaryRedirect,omitempty" toml:"sslTemporaryRedirect,omitempty" yaml:"sslTemporaryRedirect,omitempty"`
	SSLForceHost            bool              `json:"sslForceHost,omitempty" toml:"sslForceHost,omitempty" yaml:"sslForceHost,omitempty"`
	STSIncludeSubdomains    bool              `json:"stsIncludeSubdomains,omitempty" toml:"stsIncludeSubdomains,omitempty" yaml:"stsIncludeSubdomains,omitempty"`
	STSPreload              bool              `json:"stsPreload,omitempty" toml:"stsPreload,omitempty" yaml:"stsPreload,omitempty"`
	ForceSTSHeader          bool              `json:"forceSTSHeader,omitempty" toml:"forceSTSHeader,omitempty" yaml:"forceSTSHeader,omitempty"`
	FrameDeny               bool              `json:"frameDeny,omitempty" toml:"frameDeny,omitempty" yaml:"frameDeny,omitempty"`
	ContentTypeNosniff      bool              `json:"contentTypeNosniff,omitempty" toml:"contentTypeNosniff,omitempty" yaml:"contentTypeNosniff,omitempty"`
	BrowserXSSFilter        bool              `json:"browserXssFilter,omitempty" toml:"browserXssFilter,omitempty" yaml:"browserXssFilter,omitempty"`
	STSSeconds              int64             `json:"stsSeconds,omitempty" toml:"stsSeconds,omitempty" yaml:"stsSeconds,omitempty"`
	CustomBrowserXSSValue   string            `json:"customBrowserXSSValue,omitempty" toml:"customBrowserXSSValue,omitempty" yaml:"customBrowserXSSValue,omitempty"`
	CustomFrameOptionsValue string            `json:"customFrameOptionsValue,omitempty" toml:"customFrameOptionsValue,omitempty" yaml:"customFrameOptionsValue,omitempty"`
	ContentSecurityPolicy   string            `json:"contentSecurityPolicy,omitempty" toml:"contentSecurityPolicy,omitempty" yaml:"contentSecurityPolicy,omitempty"`
	PublicKey               string            `json:"publicKey,omitempty" toml:"publicKey,omitempty" yaml:"publicKey,omitempty"`
	ReferrerPolicy          string            `json:"referrerPolicy,omitempty" toml:"referrerPolicy,omitempty" yaml:"referrerPolicy,omitempty"`
	FeaturePolicy           string            `json:"featurePolicy,omitempty" toml:"featurePolicy,omitempty" yaml:"featurePolicy,omitempty"`
	SSLHost                 string            `json:"sslHost,omitempty" toml:"sslHost,omitempty" yaml:"sslHost,omitempty"`
	AllowedHosts            []string          `json:"allowedHosts,omitempty" toml:"allowedHosts,omitempty" yaml:"allowedHosts,omitempty"`
	HostsProxyHeaders       []string          `json:"hostsProxyHeaders,omitempty" toml:"hostsProxyHeaders,omitempty" yaml:"hostsProxyHeaders,omitempty"`
	SSLProxyHeaders         map[string]string `json:"sslProxyHeaders,omitempty" toml:"sslProxyHeaders,omitempty" yaml:"sslProxyHeaders,omitempty"`
}

// +k8s:deepcopy-gen=true

type Cors struct {
	AllowedOrigins     []string `json:"allowedOrigins,omitempty" toml:"allowedOrigins,omitempty" yaml:"allowedOrigins,omitempty"`
	AllowedMethods     []string `json:"allowedMethods,omitempty" toml:"allowedMethods,omitempty" yaml:"allowedMethods,omitempty"`
	AllowedHeaders     []string `json:"allowedHeaders,omitempty" toml:"allowedHeaders,omitempty" yaml:"allowedHeaders,omitempty"`
	ExposedHeaders     []string `json:"exposedHeaders,omitempty" toml:"exposedHeaders,omitempty" yaml:"exposedHeaders,omitempty"`
	MaxAge             int      `json:"maxAge,omitempty" toml:"maxAge,omitempty" yaml:"maxAge,omitempty"`
	AllowCredentials   bool     `json:"allowCredentials,omitempty" toml:"allowCredentials,omitempty" yaml:"allowCredentials,omitempty"`
	OptionsPassthrough bool
}

// +k8s:deepcopy-gen=true

type CustomHeaders struct {
	RequestHeaders  map[string]string `json:"customRequestHeaders,omitempty" toml:"customRequestHeaders,omitempty" yaml:"customRequestHeaders,omitempty"`
	ResponseHeaders map[string]string `json:"customResponseHeaders,omitempty" toml:"customResponseHeaders,omitempty" yaml:"customResponseHeaders,omitempty"`
}

// +k8s:deepcopy-gen=true

// ServersLoadBalancer holds the ServersLoadBalancer configuration.
type LoadBalancer struct {
	Sticky  *Sticky   `json:"sticky,omitempty" toml:"sticky,omitempty" yaml:"sticky,omitempty" label:"allowEmpty"`
	Servers []*Server `json:"servers,omitempty" toml:"servers,omitempty" yaml:"servers,omitempty" label-slice-as-struct:"server"`
}

// +k8s:deepcopy-gen=true

type Sticky struct {
	Cookie *Cookie `json:"cookie,omitempty" toml:"cookie,omitempty" yaml:"cookie,omitempty"`
}

// +k8s:deepcopy-gen=true

type Cookie struct {
	Name string `json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty"`

	HTTPOnly bool `json:"httpOnly,omitempty" toml:"httpOnly,omitempty" yaml:"httpOnly,omitempty"`
	Secure   bool `json:"secure,omitempty" toml:"secure,omitempty" yaml:"secure,omitempty"`

	Path   string `json:"path,omitempty" toml:"path,omitempty" yaml:"path,omitempty"`
	Domain string `json:"domain,omitempty" toml:"domain,omitempty" yaml:"domain,omitempty"`

	MaxAge   int    `json:"maxAge,omitempty" toml:"maxAge,omitempty" yaml:"maxAge,omitempty"`
	SameSite string `json:"sameSite,omitempty" toml:"sameSite,omitempty" yaml:"sameSite,omitempty"`
}

// +k8s:deepcopy-gen=true

// Server holds the server configuration.
type Server struct {
	URL    string `json:"url,omitempty" toml:"url,omitempty" yaml:"url,omitempty" label:"-"`
	Weight int    `json:"weight,omitempty" toml:"weight,omitempty" yaml:"weight,omitempty" label:"-"`
}

// +k8s:deepcopy-gen=true

type RateLimit struct {
	MaxDelay     time.Duration `json:"maxDelay,omitempty" toml:"maxDelay,omitempty" yaml:"maxDelay,omitempty" label:"-"`
	DefaultDelay time.Duration `json:"defaultDelay,omitempty" toml:"defaultDelay,omitempty" yaml:"defaultDelay,omitempty" label:"-"`
	Cooldown     time.Duration `json:"cooldown,omitempty" toml:"cooldown,omitempty" yaml:"cooldown,omitempty" label:"-"`
	Limits       []float64     `json:"limits,omitempty" toml:"limits,omitempty" yaml:"limits,omitempty" label:"-"`
}

// +k8s:deepcopy-gen=true

type HTTPCache struct {
	TTL       time.Duration `json:"ttl,omitempty" toml:"ttl,omitempty" yaml:"ttl,omitempty" label:"-"`
	KeySuffix string        `json:"key_suffix,omitempty" toml:"key_suffix,omitempty" yaml:"key_suffix,omitempty" label:"-"`
}

// +k8s:deepcopy-gen=true

type HTTPTrace struct{}

// +k8s:deepcopy-gen=true

// AccessLog holds the configuration settings for the access logger (middlewares/accesslog).
type AccessLog struct {
	Enabled       bool              `description:"Access log file enabled." json:"enabled,omitempty" toml:"enabled,omitempty" yaml:"enabled,omitempty" export:"true"`
	FilePath      string            `description:"Access log file path. Stdout is used when omitted or empty." json:"filePath,omitempty" toml:"filePath,omitempty" yaml:"filePath,omitempty" export:"true"`
	Format        string            `description:"Access log format: json | common" json:"format,omitempty" toml:"format,omitempty" yaml:"format,omitempty" export:"true"`
	Filters       *AccessLogFilters `description:"Access log filters, used to keep only specific access logs." json:"filters,omitempty" toml:"filters,omitempty" yaml:"filters,omitempty" export:"true"`
	Fields        *AccessLogFields  `description:"AccessLogFields." json:"fields,omitempty" toml:"fields,omitempty" yaml:"fields,omitempty" export:"true"`
	BufferingSize int64             `description:"Number of access log lines to process in a buffered way." json:"bufferingSize,omitempty" toml:"bufferingSize,omitempty" yaml:"bufferingSize,omitempty" export:"true"`
}

func (l *AccessLog) ToTraefikType() *traefiktypes.AccessLog {
	if l == nil {
		return nil
	}

	return &traefiktypes.AccessLog{
		FilePath:      l.FilePath,
		Format:        l.Format,
		BufferingSize: l.BufferingSize,
		Filters:       l.Filters.ToTraefikType(),
		Fields:        l.Fields.ToTraefikType(),
	}
}

// +k8s:deepcopy-gen=true

// AccessLogFilters holds filters configuration
type AccessLogFilters struct {
	StatusCodes   []string      `description:"Keep access logs with status codes in the specified range." json:"statusCodes,omitempty" toml:"statusCodes,omitempty" yaml:"statusCodes,omitempty" export:"true"`
	RetryAttempts bool          `description:"Keep access logs when at least one retry happened." json:"retryAttempts,omitempty" toml:"retryAttempts,omitempty" yaml:"retryAttempts,omitempty" export:"true"`
	MinDuration   time.Duration `description:"Keep access logs when request took longer than the specified duration." json:"minDuration,omitempty" toml:"minDuration,omitempty" yaml:"minDuration,omitempty" export:"true"`
}

func (f *AccessLogFilters) ToTraefikType() *traefiktypes.AccessLogFilters {
	if f == nil {
		return nil
	}

	return &traefiktypes.AccessLogFilters{
		StatusCodes:   f.StatusCodes,
		RetryAttempts: f.RetryAttempts,
		MinDuration:   traefiktypes2.Duration(f.MinDuration),
	}
}

// +k8s:deepcopy-gen=true

// AccessLogFields holds configuration for access log fields
type AccessLogFields struct {
	DefaultMode string            `description:"Default mode for fields: keep | drop" json:"defaultMode,omitempty" toml:"defaultMode,omitempty" yaml:"defaultMode,omitempty"  export:"true"`
	Names       map[string]string `description:"Override mode for fields" json:"names,omitempty" toml:"names,omitempty" yaml:"names,omitempty" export:"true"`
	Headers     *FieldHeaders     `description:"Headers to keep, drop or redact" json:"headers,omitempty" toml:"headers,omitempty" yaml:"headers,omitempty" export:"true"`
}

func (f *AccessLogFields) ToTraefikType() *traefiktypes.AccessLogFields {
	if f == nil {
		return nil
	}

	return &traefiktypes.AccessLogFields{
		DefaultMode: f.DefaultMode,
		Names:       f.Names,
		Headers:     f.Headers.ToTraefikType(),
	}
}

// +k8s:deepcopy-gen=true

// FieldHeaders holds configuration for access log headers
type FieldHeaders struct {
	DefaultMode string            `description:"Default mode for fields: keep | drop | redact" json:"defaultMode,omitempty" toml:"defaultMode,omitempty" yaml:"defaultMode,omitempty" export:"true"`
	Names       map[string]string `description:"Override mode for headers" json:"names,omitempty" toml:"names,omitempty" yaml:"names,omitempty" export:"true"`
}

func (h *FieldHeaders) ToTraefikType() *traefiktypes.FieldHeaders {
	if h == nil {
		return nil
	}

	return &traefiktypes.FieldHeaders{
		DefaultMode: h.DefaultMode,
		Names:       h.Names,
	}
}
