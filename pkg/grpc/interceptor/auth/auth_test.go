package grpcauth

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"google.golang.org/grpc/metadata"
)

var (
	idToken                           = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiaHR0cHM6Ly9hdXRoMC5jb20vYXBpL3YyLyJdLCJleHAiOjE1NzkxNjc0MTQsImh0dHA6Ly9vcmNoZXN0cmF0ZS5pbmZvIjp7InRlbmFudF9pZCI6ImI0OWVlMWJjLWYwZmEtNDMwZC04OWIyLWE0ZmQwZGM5ODkwNiIsInJvbGUiOiJ0ZXN0LXJvbGUifSwiaWF0IjoxNTc5MTYzODE0LCJpc3MiOiJPcmNoZXN0cmF0ZSIsImp0aSI6IjZlZmY3MzI0LTVkZTEtNDA2NS05NGNmLWU3ZWYzZTliYjg1MCIsIm5iZiI6MTU3OTE2MzgxNCwic2NwIjpbInJlYWQ6dXNlcnMiLCJ1cGRhdGU6dXNlcnMiLCJjcmVhdGU6dXNlcnMiXSwic3ViIjoiZTJlLXRlc3QifQ.fvlJcrCwbvj-W1VrfSzcn5F7LpsZ0xbOQTcCqVwwmyq8EOv5VwoV-geoX6tj4d0T2pew-6EK8DR-GrwXVjlo2LQQhYY_TRpnVHl1wDE1IvahExnh_0oPwpH3oKjsxbLPyM94bG-eIJGyInA3w-llCXR5WhOwccO4lKW4GaAXsj6TKGiowh_9HEw9jSN2y9OXGvUiE9_8n_5rp1Shp_vBMHJ-5usOozoaJdgl13Dln1YTqSl422CKb1UndBGRXayCfMpqnzLuURTYYspWOn3c6QTbjjMAm8ifZIl8rDrI8zl8j2FM1kHZt-5ZZe5zJv7rCGwPQviLnWQBqIVElJv6Tg"
	idTokenExpired                    = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6Ik1UTXpRamRETUROR01qazVOREpCTUVKQk56ZzBSVUk0UWpjek5UZzJNMFUxT0VZNVJrUTNRdyJ9.eyJodHRwOi8vdGVuYW50LmluZm8vdGVuYW50X2lkIjoiMTkwZTBlMmItMmZiNS00NGEwLTllNDgtYzUyYWM0Mzg0MzI5IiwiaHR0cDovL3RlbmFudC5pbmZvL3RlbmFudF9yb2xlIjoidXNlciIsImh0dHA6Ly90ZW5hbnQuaW5mby90ZW5hbnRfY29tcGFnbnkiOiJDb2RlRmkiLCJuaWNrbmFtZSI6ImJhciIsIm5hbWUiOiJiYXJAZm9vLmNvbSIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9kYzhhNDJhYmEzNjUxYjBiMWYwODhlZjkyOGZmM2IxZD9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmJhLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDE5LTExLTI4VDA5OjM3OjI5LjkxN1oiLCJpc3MiOiJodHRwczovL2Rldi1iZDZlM2psYy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NWRkZThiMWI4YjU5YjEwZTE5ODU0ODEzIiwiYXVkIjoiWkNlN0p1Q2xpdTIwUjBzTHBTRTB3OEk3V3dhMTYyV2QiLCJpYXQiOjE1NzQ5MzM4NTEsImV4cCI6MTU3NDk2OTg1MX0.advUv8dSHnF2Tj0NAO-hFMJD-H0Y55FbxaOM_x-qZWNTKo1ycdfVy3-i1ODJgmdyLNrJhKpOMuEEg61eqsULG5Fre79bmErHI9UEmKLeY1fcfboR1J9vxgiyNcBtoV4F2CzpXWo-Xp_-Fhkam2jJ-GwdY3wRT9IM4GikJosZqzbhieqm44irhHp3O-afAhU-5xm4eybz1FP67_t8xHPnGIoIQlxUXeKN8AwjmWMIoe6mdlHYyoFAtt05hL48XvmH-IvOVXn7bi3CBytnBm_FudWtdnyddW-TSZ9IhhFR7zWm4Tsg3NPRVqtG6HvONwtiaz-IArcd-RsVDascx_tO1g"
	accessTokenUntrustedSigner        = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IlJFRTVSRUpETURFeFJrVTFNVGhDUTBGQ05rTXdSVEkyTVVSRk1qQXlOekUyTjBNMU1rWXpNZyJ9.eyJodHRwczovL2Rldi5hYmNjZC5jb20iOnsicm9sZXMiOlsiVHJhZGUgRXhlY3V0b3IiXSwibGVnYWxFbnRpdGllcyI6WzFdLCJsZWdhbEVudGl0eSI6MX0sImlzcyI6Imh0dHBzOi8vZGV2LWFiY2NkLmV1LmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1ZDhjODFmMTRiMmZlYjBkNzA0YzVlMGYiLCJhdWQiOlsiYmVhdC1hdXRoMC1hcGkiLCJodHRwczovL2Rldi1hYmNjZC5ldS5hdXRoMC5jb20vdXNlcmluZm8iXSwiaWF0IjoxNTc1NTQxMjA1LCJleHAiOjE1NzU2Mjc2MDUsImF6cCI6ImVkQ3hvVlRPcFhCN2t0NmdpVFdoOG1CZ0pnTVdvTzJ2Iiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInBlcm1pc3Npb25zIjpbImNhOmNyZWF0ZSIsImNhOnVwZGF0ZSIsImNhOnZpZXdEZXRhaWxzIiwibGVnYWwtZW50aXR5OnJlYWQiLCJtYXN0ZXItZGF0YTpyZWFkIiwibm9taW5hdGlvbjpjcmVhdGUiLCJub21pbmF0aW9uOnVwZGF0ZSIsIm5vbWluYXRpb246dmlld0RldGFpbHMiLCJ2b3lhZ2U6Y3JlYXRlIiwidm95YWdlOnJlYWQiXX0.jJnJjTLHsElFU3O7xKuh7jL1ho9-Z7Jxco16hDxoRg_TFdOCN82wVeJHZbDjdLqjV0k4F05YWEmFWn7CEAmr43ndoprsAr3OfBnrjYKyJ4oqiguPAUakBqoLtaEE-AsxyQmCzZGwKXHtNMDIhh0vwHVASdHTwxiApumRWfEXzmu5pmOYwoTJ8vVSUVCCDG3hL6u4UxYdng30XlWgbn_Szlaq9sIoIllZOL8vn4hkkW98CQfjexpaYDjywVfbPD3-TSSznHiF6TvmogCttkb73hbJF246hq-guR0nfdQm1ivAUkzXcUOql6QtHvYgdrzw5xPOqNIMihFvIK8XRCZ_pw"
	accessTokenWithoutTenantID        = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiaHR0cHM6Ly9hdXRoMC5jb20vYXBpL3YyLyJdLCJleHAiOjE1NzkxNjc0NjksImh0dHA6Ly9vcmNoZXN0cmF0ZS5pbmZvIjp7InRlbmFudF9pZCI6IiIsInJvbGUiOiJ0ZXN0LXJvbGUifSwiaWF0IjoxNTc5MTYzODY5LCJpc3MiOiJPcmNoZXN0cmF0ZSIsImp0aSI6IjNkNDAyYWFlLTMwY2YtNDcxNy04MGRmLTg4ODE4OTJhOTUxOCIsIm5iZiI6MTU3OTE2Mzg2OSwic2NwIjpbInJlYWQ6dXNlcnMiLCJ1cGRhdGU6dXNlcnMiLCJjcmVhdGU6dXNlcnMiXSwic3ViIjoiZTJlLXRlc3QifQ.nTr2eY8mXXD6kqUnhx5pAydwUXnxpzPdZ9qqMcPaDEsNSJT_HJvYc11kut7VN_DVL3sFT6xo1auB40w96xh1TatYGYB3FmISfIbZ4XAjgkRzTB5uaf8eoi0DDnAQ3ycxVdmuKDapVW5gS9FQmoGOwcC_ojoQtQKUc3XyTiHAowurTKSre329EunCEj2dMSRBTEmg_vnWgGmgtpxOI9f4l1hrrQ3FAGbZobdVoqkTzLwVqo1GblxUioQGYPSy6okO6XPKL2G0P62iIJqClhNRQP0pZJHucJCipdZYaOLrBdepO7srIUt4gM3qXkkWohPDqOujdUoqUUtSq6C37rwGtA"
	certificateOneLineOrchestrateTest = "MIIDYjCCAkoCCQC9pJWk7qdipjANBgkqhkiG9w0BAQsFADBzMQswCQYDVQQGEwJGUjEOMAwGA1UEBwwFUGFyaXMxEjAQBgNVBAoMCUNvbnNlblN5czEQMA4GA1UECwwHUGVnYVN5czEuMCwGA1UEAwwlZTJlLXRlc3RzLm9yY2hlc3RyYXRlLmNvbnNlbnN5cy5wYXJpczAeFw0xOTEyMjcxNjI5MTdaFw0yMDEyMjYxNjI5MTdaMHMxCzAJBgNVBAYTAkZSMQ4wDAYDVQQHDAVQYXJpczESMBAGA1UECgwJQ29uc2VuU3lzMRAwDgYDVQQLDAdQZWdhU3lzMS4wLAYDVQQDDCVlMmUtdGVzdHMub3JjaGVzdHJhdGUuY29uc2Vuc3lzLnBhcmlzMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAo0NqWqI3TSi1uOBvCUquclWo4LcsYT21tNUXQ8YyqVYRSsiBv+ZKZBCjD8XklLPih40kFSe+r6DNca5/LH/okQIdc8nsQg+BLCkXeH2NFv+QYtPczAw4YhS6GVxJk3u9sfp8NavWBcQbD3MMDpehMOvhSl0zoP/ZlH6ErKHNtoQgUpPNVQGysNU21KpClmIDD/L1drsbq+rFiDrcVWaOLwGxr8SBd/0b4ngtcwH16RJaxcIXXT5AVia1CNdzmU5/AIg3OfgzvKn5AGrMZBsmGAiCyn4/P3PnuF81/WHukk5ETLnzOH+vC2elSmZ8y80HCGeqOiQ1rs66L936wX8cDwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQCNcTs3n/Ps+yIZDH7utxTOaqpDTCB10MzPmb22UAal89couIT6R0fAu14p/LTkxdb2STDySsQY2/Lv6rPdFToHGUI9ZYOTYW1GOWkt1EAao9BzdsoJVwmTON6QnOBKy/9RxlhWP+XSWVsY0te6KYzS7rQyzQoJQeeBNMpUnjiQji9kKi5j9rbVMdjIb4HlmYrcE95ps+oFkyJoA1HLVytAeOjJPXGToNlv3k2UPJzOFUM0ujWWeBTyHMCmZ4RhlrfzDNffY5dlW82USjc5dBlzRyZalXSjhcVhK4asUodomVntrvCShp/8C9LpbQZ+ugFNE8J6neStWrhpRU9/sBJx"
)

// TODO: adding tests with APIKey Header and TenantID header

func TestAuth(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *authjwt.Config
		token    string
		want     string
		errCode  uint64
		isValide bool
	}{
		{
			"nominal case",
			&authjwt.Config{
				SkipClaimsValidation: true,
				Certificate:          []byte(certificateOneLineOrchestrateTest),
				ClaimsNamespace:      "http://orchestrate.info",
			},
			idToken,
			"b49ee1bc-f0fa-430d-89b2-a4fd0dc98906",
			0,
			true,
		},
		{
			"expired filed",
			&authjwt.Config{
				Certificate:     []byte(certificateOneLineOrchestrateTest),
				ClaimsNamespace: "http://orchestrate.info",
			},
			idTokenExpired,
			"",
			errors.Unauthorized,
			false,
		},
		{
			"empty bearer",
			&authjwt.Config{
				Certificate:     []byte(certificateOneLineOrchestrateTest),
				ClaimsNamespace: "http://orchestrate.info",
			},
			"",
			"",
			errors.Unauthorized,
			false,
		},
		{
			"Untrusted Signer",
			&authjwt.Config{
				Certificate:          []byte(certificateOneLineOrchestrateTest),
				ClaimsNamespace:      "http://orchestrate.info",
				SkipClaimsValidation: true,
			},
			accessTokenUntrustedSigner,
			"",
			errors.Unauthorized,
			false,
		},
		{
			"Token without TenantID",
			&authjwt.Config{
				Certificate:          []byte(certificateOneLineOrchestrateTest),
				ClaimsNamespace:      "http://orchestrate.info",
				SkipClaimsValidation: true,
			},
			accessTokenWithoutTenantID,
			"",
			errors.PermissionDenied,
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := setupContext(test.token)
			checker, err := authjwt.New(test.cfg)
			require.NoError(t, err)
			authFunc := Auth(checker, true)
			checkedCtx, err := authFunc(ctx)
			if !test.isValide {
				assert.Error(t, err, "Auth should error")
				assert.Equal(t, test.errCode, errors.FromError(err).Code, "Error should be correct")
			} else {
				assert.NoError(t, err, "Auth should not error")
				token := authjwt.FromContext(checkedCtx)
				assert.NotNil(t, token, "Token should not be nil")
				tenantID := token.Claims.(*authjwt.Claims).Orchestrate.TenantID
				assert.Equal(t, test.want, tenantID, "Tenant UUID should be correct")
			}
		})
	}
}

func setupContext(tokenValue string) context.Context {
	if tokenValue == "" {
		return context.TODO()
	}

	md := metadata.Pairs("Authorization", tokenValue)
	ctx := metautils.NiceMD(md).ToIncoming(context.TODO())

	return ctx
}
