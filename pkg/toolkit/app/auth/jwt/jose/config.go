package jose

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(issuerURLViperKey, issuerURLEnv)
	_ = viper.BindEnv(audienceViperKey, audienceEnv)
	_ = viper.BindEnv(orchestrateClaimsViperKey, orchestrateClaimsEnv)
}

func Flags(f *pflag.FlagSet) {
	issuerURLFlags(f)
	audienceFlags(f)
	orchestrateClaimPath(f)
}

const (
	issuerURLFlag     = "auth-jwt-issuer-url"
	issuerURLViperKey = "auth.jwt.issuer-url"
	issuerURLEnv      = "AUTH_JWT_ISSUER_URL"
)

const (
	audienceFlag     = "auth-jwt-audience"
	audienceViperKey = "auth.jwt.audience"
	audienceEnv      = "AUTH_JWT_AUDIENCE"
)

const (
	orchestrateClaimsFlag     = "auth-jwt-orchestrate-claims"
	orchestrateClaimsViperKey = "auth.jwt.orchestrate.claims"
	orchestrateClaimsEnv      = "AUTH_JWT_ORCHESTRATE_CLAIMS"
)

func issuerURLFlags(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`JWT issuer server domain (ie. https://orchestrate.eu.auth0.com).
Environment variable: %q`, issuerURLEnv)
	f.String(issuerURLFlag, "", desc)
	_ = viper.BindPFlag(issuerURLViperKey, f.Lookup(issuerURLFlag))
}

func audienceFlags(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Expected audience ("aud" field) of JWT tokens.
Environment variable: %q`, audienceEnv)
	f.StringSlice(audienceFlag, []string{}, desc)
	_ = viper.BindPFlag(audienceViperKey, f.Lookup(audienceFlag))
}

func orchestrateClaimPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Path to for orchestrate claims in the Access Token.
Environment variable: %q`, orchestrateClaimsEnv)
	f.String(orchestrateClaimsFlag, "", desc)
	_ = viper.BindPFlag(orchestrateClaimsViperKey, f.Lookup(orchestrateClaimsFlag))
}

type Config struct {
	IssuerURL         string
	CacheTTL          time.Duration
	Audience          []string
	OrchestrateClaims string
}

func NewConfig(vipr *viper.Viper) *Config {
	issuerURL := vipr.GetString(issuerURLViperKey)

	if issuerURL == "" {
		return nil
	}

	cfg := &Config{
		IssuerURL: issuerURL,
		Audience:  vipr.GetStringSlice(audienceViperKey),
		CacheTTL:  5 * time.Minute, // TODO: Make the cache ttl an ENV var if needed
	}

	if cPath := vipr.GetString(orchestrateClaimsViperKey); cPath != "" {
		cfg.OrchestrateClaims = cPath
	}

	return cfg
}
