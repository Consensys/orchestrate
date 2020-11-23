package jwt_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/jwt"
	jwtgenerator "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/jwt/generator"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tls/certificate"
	tlstestutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/tls/testutils"
)

func TestJWT(t *testing.T) {
	tests := []struct { //nolint:maligned // reason
		name                       string
		cfg                        *jwt.Config
		genCfg                     *jwtgenerator.Config
		jwtTenantID                string
		jwtTTL                     time.Duration
		tenantID                   string
		expectedValidAuth          bool
		expectedImpersonatedTenant string
		expectedAllowedTenants     []string
	}{
		{
			"invalid signature",
			&jwt.Config{
				ClaimsNamespace: "orchestrate.test",
				Certificate:     []byte(tlstestutils.OneLineRSACertPEMB),
			},
			&jwtgenerator.Config{
				KeyPair: &certificate.KeyPair{
					Cert: []byte(tlstestutils.OneLineRSACertPEMA),
					Key:  []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
				ClaimsNamespace: "orchestrate.test",
			},
			"foo",
			10 * time.Hour,
			"foo",
			false,
			"",
			nil,
		},
		{
			"expired token",
			&jwt.Config{
				ClaimsNamespace: "orchestrate.test",
				Certificate:     []byte(tlstestutils.OneLineRSACertPEMA),
			},
			&jwtgenerator.Config{
				KeyPair: &certificate.KeyPair{
					Cert: []byte(tlstestutils.OneLineRSACertPEMA),
					Key:  []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
				ClaimsNamespace: "orchestrate.test",
			},
			"foo",
			-10 * time.Second,
			"foo",
			false,
			"",
			nil,
		},
		{
			"distinct jwt custom claims namespace",
			&jwt.Config{
				ClaimsNamespace: "orchestrate.prod",
				Certificate:     []byte(tlstestutils.OneLineRSACertPEMA),
			},
			&jwtgenerator.Config{
				KeyPair: &certificate.KeyPair{
					Cert: []byte(tlstestutils.OneLineRSACertPEMA),
					Key:  []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
				ClaimsNamespace: "orchestrate.dev",
			},
			"foo",
			10 * time.Hour,
			"foo",
			false,
			"",
			nil,
		},
		{
			"empty tenant id",
			&jwt.Config{
				ClaimsNamespace: "orchestrate.test",
				Certificate:     []byte(tlstestutils.OneLineRSACertPEMA),
			},
			&jwtgenerator.Config{
				KeyPair: &certificate.KeyPair{
					Cert: []byte(tlstestutils.OneLineRSACertPEMA),
					Key:  []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
				ClaimsNamespace: "orchestrate.test",
			},
			"",
			10 * time.Hour,
			"",
			false,
			"",
			nil,
		},
		{
			"JWT foo accessing empty tenant",
			&jwt.Config{
				ClaimsNamespace: "orchestrate.test",
				Certificate:     []byte(tlstestutils.OneLineRSACertPEMA),
			},
			&jwtgenerator.Config{
				KeyPair: &certificate.KeyPair{
					Cert: []byte(tlstestutils.OneLineRSACertPEMA),
					Key:  []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
				ClaimsNamespace: "orchestrate.test",
			},
			"foo",
			10 * time.Hour,
			"",
			true,
			"foo",
			[]string{"foo", multitenancy.DefaultTenant},
		},
		{
			"JWT foo accessing foo tenant",
			&jwt.Config{
				ClaimsNamespace: "orchestrate.test",
				Certificate:     []byte(tlstestutils.OneLineRSACertPEMA),
			},
			&jwtgenerator.Config{
				KeyPair: &certificate.KeyPair{
					Cert: []byte(tlstestutils.OneLineRSACertPEMA),
					Key:  []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
				ClaimsNamespace: "orchestrate.test",
			},
			"foo",
			10 * time.Hour,
			"foo",
			true,
			"foo",
			[]string{"foo"},
		},
		{
			"JWT foo accessing bar tenant",
			&jwt.Config{
				ClaimsNamespace: "orchestrate.test",
				Certificate:     []byte(tlstestutils.OneLineRSACertPEMA),
			},
			&jwtgenerator.Config{
				KeyPair: &certificate.KeyPair{
					Cert: []byte(tlstestutils.OneLineRSACertPEMA),
					Key:  []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
				ClaimsNamespace: "orchestrate.test",
			},
			"foo",
			10 * time.Hour,
			"bar",
			false,
			"",
			nil,
		},
		{
			"JWT * accessing empty tenant",
			&jwt.Config{
				ClaimsNamespace: "orchestrate.test",
				Certificate:     []byte(tlstestutils.OneLineRSACertPEMA),
			},
			&jwtgenerator.Config{
				KeyPair: &certificate.KeyPair{
					Cert: []byte(tlstestutils.OneLineRSACertPEMA),
					Key:  []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
				ClaimsNamespace: "orchestrate.test",
			},
			"*",
			10 * time.Hour,
			"",
			true,
			"_",
			[]string{multitenancy.Wildcard},
		},
		{
			"JWT * accessing foo tenant",
			&jwt.Config{
				ClaimsNamespace: "orchestrate.test",
				Certificate:     []byte(tlstestutils.OneLineRSACertPEMA),
			},
			&jwtgenerator.Config{
				KeyPair: &certificate.KeyPair{
					Cert: []byte(tlstestutils.OneLineRSACertPEMA),
					Key:  []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
				ClaimsNamespace: "orchestrate.test",
			},
			"*",
			10 * time.Hour,
			"foo",
			true,
			"foo",
			[]string{"foo", multitenancy.DefaultTenant},
		},
		{
			"JWT * accessing default tenant",
			&jwt.Config{
				ClaimsNamespace: "orchestrate.test",
				Certificate:     []byte(tlstestutils.OneLineRSACertPEMA),
			},
			&jwtgenerator.Config{
				KeyPair: &certificate.KeyPair{
					Cert: []byte(tlstestutils.OneLineRSACertPEMA),
					Key:  []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
				ClaimsNamespace: "orchestrate.test",
			},
			"*",
			10 * time.Hour,
			multitenancy.DefaultTenant,
			true,
			multitenancy.DefaultTenant,
			[]string{multitenancy.DefaultTenant},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker, err := jwt.New(tt.cfg)
			require.NoError(t, err)

			gen, err := jwtgenerator.New(tt.genCfg)
			require.NoError(t, err)

			token, err := gen.GenerateAccessTokenWithTenantID(tt.jwtTenantID, tt.jwtTTL)
			require.NoError(t, err)

			ctx := authutils.WithAuthorization(context.Background(), fmt.Sprintf("Bearer %v", token))
			ctx = multitenancy.WithTenantID(ctx, tt.tenantID)

			checkedCtx, err := checker.Check(ctx)
			if tt.expectedValidAuth {
				require.NoError(t, err, "Authentication check should succeeds")
				assert.Equal(t, tt.expectedImpersonatedTenant, multitenancy.TenantIDFromContext(checkedCtx))
				assert.Equal(t, tt.expectedAllowedTenants, multitenancy.AllowedTenantsFromContext(checkedCtx))
			} else {
				require.Error(t, err, "Authentication should fail")
			}
		})
	}
}
