package jwt

import (
	"context"
	"fmt"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
)

// TODO: adding new tests to add coverage

const (
	idTokenExpired                    = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6Ik1UTXpRamRETUROR01qazVOREpCTUVKQk56ZzBSVUk0UWpjek5UZzJNMFUxT0VZNVJrUTNRdyJ9.eyJodHRwOi8vdGVuYW50LmluZm8vdGVuYW50X2lkIjoiMTkwZTBlMmItMmZiNS00NGEwLTllNDgtYzUyYWM0Mzg0MzI5IiwiaHR0cDovL3RlbmFudC5pbmZvL3RlbmFudF9yb2xlIjoidXNlciIsImh0dHA6Ly90ZW5hbnQuaW5mby90ZW5hbnRfY29tcGFnbnkiOiJDb2RlRmkiLCJuaWNrbmFtZSI6ImJhciIsIm5hbWUiOiJiYXJAZm9vLmNvbSIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci9kYzhhNDJhYmEzNjUxYjBiMWYwODhlZjkyOGZmM2IxZD9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmJhLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDE5LTExLTI4VDA5OjM3OjI5LjkxN1oiLCJpc3MiOiJodHRwczovL2Rldi1iZDZlM2psYy5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NWRkZThiMWI4YjU5YjEwZTE5ODU0ODEzIiwiYXVkIjoiWkNlN0p1Q2xpdTIwUjBzTHBTRTB3OEk3V3dhMTYyV2QiLCJpYXQiOjE1NzQ5MzM4NTEsImV4cCI6MTU3NDk2OTg1MX0.advUv8dSHnF2Tj0NAO-hFMJD-H0Y55FbxaOM_x-qZWNTKo1ycdfVy3-i1ODJgmdyLNrJhKpOMuEEg61eqsULG5Fre79bmErHI9UEmKLeY1fcfboR1J9vxgiyNcBtoV4F2CzpXWo-Xp_-Fhkam2jJ-GwdY3wRT9IM4GikJosZqzbhieqm44irhHp3O-afAhU-5xm4eybz1FP67_t8xHPnGIoIQlxUXeKN8AwjmWMIoe6mdlHYyoFAtt05hL48XvmH-IvOVXn7bi3CBytnBm_FudWtdnyddW-TSZ9IhhFR7zWm4Tsg3NPRVqtG6HvONwtiaz-IArcd-RsVDascx_tO1g"
	idToken                           = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6Ik1UTXpRamRETUROR01qazVOREpCTUVKQk56ZzBSVUk0UWpjek5UZzJNMFUxT0VZNVJrUTNRdyJ9.eyJodHRwOi8vdGVuYW50LmluZm8vdGVuYW50X2lkIjoiYjQ5ZWUxYmMtZjBmYS00MzBkLTg5YjItYTRmZDBkYzk4OTA2IiwiaHR0cDovL3RlbmFudC5pbmZvL3RlbmFudF9yb2xlIjoidXNlciIsImh0dHA6Ly90ZW5hbnQuaW5mby90ZW5hbnRfY29tcGFnbnkiOiJQZWdhU3lzIiwibmlja25hbWUiOiJmb28iLCJuYW1lIjoiZm9vQGJhci5jb20iLCJwaWN0dXJlIjoiaHR0cHM6Ly9zLmdyYXZhdGFyLmNvbS9hdmF0YXIvZjNhZGE0MDVjZTg5MGI2ZjgyMDQwOTRkZWIxMmQ4YTg_cz00ODAmcj1wZyZkPWh0dHBzJTNBJTJGJTJGY2RuLmF1dGgwLmNvbSUyRmF2YXRhcnMlMkZmby5wbmciLCJ1cGRhdGVkX2F0IjoiMjAxOS0xMi0wNlQwOTo0ODowMS41NTNaIiwiaXNzIjoiaHR0cHM6Ly9kZXYtYmQ2ZTNqbGMuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfDVkZGU4ZTYyNzY5YTJkMGVkM2FmNTM4ZSIsImF1ZCI6IlpDZTdKdUNsaXUyMFIwc0xwU0UwdzhJN1d3YTE2MldkIiwiaWF0IjoxNTc1NjI1NjgyLCJleHAiOjE1NzU3MTIwODJ9.muHMxGe0EaSYnRCVpVAPeIfeEr4VLnN54DcWOxk6CMBUlNq2gzElxiKkZ2IUS6oZXCwHvob40mMJQJyIPpRBn23ZsIZLK3Iy4Xbf-TytvtSKWMX4Jiw1WgNey7_DsjHtT6Wi9OufS2NF49sK39m0hDXf2GCqqtYFg5XNQLMujfDdplxN2gRHP3VEey3PtSMBFIdlAkv2mCA5SPBlxmkCtGmgiQa223bPl2rnCA5PF7XjNVTg2v59m34ADZ8cR-J6h1UrKPXFmCXEO1gHC_wpiN7E0pjjnJVORDN27b5zAASADPSh9tyZlWbZa14SAP8M9gzOChS5z5b31efuvA8Rxw"
	idTokenNominalUsecase             = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiaHR0cHM6Ly9hdXRoMC5jb20vYXBpL3YyLyJdLCJleHAiOjE1ODAyMjMwNTcsImh0dHA6Ly9vcmNoZXN0cmF0ZS5pbmZvIjp7InRlbmFudF9pZCI6ImYzMGM0NTJiLWU1ZmItNDEwMi1hNDVkLWJjMDBhMDYwYmNjNiJ9LCJpYXQiOjE1ODAxMzY2NTcsImlzcyI6Ik9yY2hlc3RyYXRlIiwianRpIjoiZTFiMjE2NDYtZjUxYS00MTk1LWE4OTYtYjJiOGMxNzQyYTExIiwibmJmIjoxNTgwMTM2NjU3LCJzY3AiOlsicmVhZDp1c2VycyIsInVwZGF0ZTp1c2VycyIsImNyZWF0ZTp1c2VycyJdLCJzdWIiOiJlMmUtdGVzdCJ9.d_V25rftcW1nV7LLyRuug_U9AQRNghLQr4QB3_a6CBe5BazXZNer4sXIR2WZEFzJVFHkeQjPihJ96eHee3Kwf7GL-ElZ9uGag45pKjjiIJaI4fgne8PcAGKYqmqNAL7_oi7_KJ_PLTIUQDWcKzI0MegBZ6-WTY5qoXDsq-HoVDLMN0MO2lC3gkApoOv-TLTDQapPBeAT_wOc18L2juAy_9PeZz99axkrYN3w68Rbf1irSzMieXpmjcQOxo5qdoP0KQ9zZztCfsSqyEgOXAfgsslSG3E__AxYzXR285Mp_lGdIVEoVlmhQDvhWz3og8lfNj2f_-bX1yoq4NQddxfGHQ"
	certificateSuccessOrchestrateTest = `
	MIIDYjCCAkoCCQC9pJWk7qdipjANBgkqhkiG9w0BAQsFADBzMQswCQYDVQQGEwJG
	UjEOMAwGA1UEBwwFUGFyaXMxEjAQBgNVBAoMCUNvbnNlblN5czEQMA4GA1UECwwH
	UGVnYVN5czEuMCwGA1UEAwwlZTJlLXRlc3RzLm9yY2hlc3RyYXRlLmNvbnNlbnN5
	cy5wYXJpczAeFw0xOTEyMjcxNjI5MTdaFw0yMDEyMjYxNjI5MTdaMHMxCzAJBgNV
	BAYTAkZSMQ4wDAYDVQQHDAVQYXJpczESMBAGA1UECgwJQ29uc2VuU3lzMRAwDgYD
	VQQLDAdQZWdhU3lzMS4wLAYDVQQDDCVlMmUtdGVzdHMub3JjaGVzdHJhdGUuY29u
	c2Vuc3lzLnBhcmlzMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAo0Nq
	WqI3TSi1uOBvCUquclWo4LcsYT21tNUXQ8YyqVYRSsiBv+ZKZBCjD8XklLPih40k
	FSe+r6DNca5/LH/okQIdc8nsQg+BLCkXeH2NFv+QYtPczAw4YhS6GVxJk3u9sfp8
	NavWBcQbD3MMDpehMOvhSl0zoP/ZlH6ErKHNtoQgUpPNVQGysNU21KpClmIDD/L1
	drsbq+rFiDrcVWaOLwGxr8SBd/0b4ngtcwH16RJaxcIXXT5AVia1CNdzmU5/AIg3
	OfgzvKn5AGrMZBsmGAiCyn4/P3PnuF81/WHukk5ETLnzOH+vC2elSmZ8y80HCGeq
	OiQ1rs66L936wX8cDwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQCNcTs3n/Ps+yIZ
	DH7utxTOaqpDTCB10MzPmb22UAal89couIT6R0fAu14p/LTkxdb2STDySsQY2/Lv
	6rPdFToHGUI9ZYOTYW1GOWkt1EAao9BzdsoJVwmTON6QnOBKy/9RxlhWP+XSWVsY
	0te6KYzS7rQyzQoJQeeBNMpUnjiQji9kKi5j9rbVMdjIb4HlmYrcE95ps+oFkyJo
	A1HLVytAeOjJPXGToNlv3k2UPJzOFUM0ujWWeBTyHMCmZ4RhlrfzDNffY5dlW82U
	Sjc5dBlzRyZalXSjhcVhK4asUodomVntrvCShp/8C9LpbQZ+ugFNE8J6neStWrhp
	RU9/sBJx`
	certificateExpectedClientTestEnv = `-----BEGIN CERTIFICATE-----
MIIDBzCCAe+gAwIBAgIJBCTenp/s9+rWMA0GCSqGSIb3DQEBCwUAMCExHzAdBgNV
BAMTFmRldi1hYmNjZC5ldS5hdXRoMC5jb20wHhcNMTkwOTE1MTQyMjU2WhcNMzMw
NTI0MTQyMjU2WjAhMR8wHQYDVQQDExZkZXYtYWJjY2QuZXUuYXV0aDAuY29tMIIB
IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqr+P/k4219FM9vlartkxQZSD
WaWkoEUZxBsb0E2udxF0YM1JaCZbBP5YjKOVq3vVp85xuviDMFaeSEVDE4VT4uqq
llHZnl1NNHAAnMkzyOQ6BLTMpZPoxIYY2BVPjo1wePZXCW09q8EHXXxlOM0Ba8A6
mAjLnh8VZhnHFBfVfuokySvVlHS7xpKn9AWkit5ffrd78wHJMnPvRhocPyUR/JHN
OY3JojVb3wr12sOOqDwXdnRrL1PYbWOj4TcyfD5wgYIKKJKkp16R75GPllmOkN3X
MOxnrag5WDpWGKfXJ92hC9FVEkqCNxVB9B+GZjUwpH7OdZ+/ofuPaqElxqtEEQID
AQABo0IwQDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSRxTaA/4JQWQZjxZS4
iRYsk/dUwTAOBgNVHQ8BAf8EBAMCAoQwDQYJKoZIhvcNAQELBQADggEBACHYUdwR
hCKnm3e0nSqTtZbijNvFX4ddHMfNILO1GdwDkxjvSistFCtBmUy5EFQoHTyS6iZo
hLEgXwdo2RC3ZjS3hUr6JPamJhe5rW/w1DEJP2t4iBtXrxe2injGzXMh1UbNslWN
cwrPZvDB8nD5g320arFkI7M+tyDNUPOCzPa/b9D76rmHdzP9BkXvVmGGrS3Ie1RM
dhT9e5c4fcqmBv02p+eyPwJMpjzy8owqNyYzR9JZDhfX9C57hALXUmoYjXPJg2U3
/qUhcZpjmBihcd3bdLhb+8SK6RweVu4dN3/diqcTZsSQQDo61ThEV6M1ktLDmb7x
gvimsav4koWYWME=
-----END CERTIFICATE-----`
)

func TestJWT(t *testing.T) {
	tests := []struct { //nolint:maligned // reason
		name             string
		cfg              *Config
		token            string
		signatureIsValid bool
		errValue         uint32
		isInvalid        bool
	}{
		{
			"nominal usecase",
			&Config{
				SkipClaimsValidation: true,
				ClaimsNamespace:      "http://orchestrate.info",
				Certificate:          []byte(certificateSuccessOrchestrateTest),
			},
			idTokenNominalUsecase,
			true,
			0,
			false,
		},
		{
			"signature is not valid",
			&Config{
				SkipClaimsValidation: true,
				ClaimsNamespace:      "http://orchestrate.info",
				Certificate:          []byte(certificateExpectedClientTestEnv),
			},
			idToken,
			false,
			jwt.ValidationErrorSignatureInvalid,
			true,
		},
		{
			"expired filed",
			&Config{
				ClaimsNamespace: "http://orchestrate.info",
				Certificate:     []byte(certificateSuccessOrchestrateTest),
			},
			idTokenExpired,
			false,
			jwt.ValidationErrorExpired,
			true,
		},
		{
			"no tenant id",
			&Config{
				SkipClaimsValidation: true,
				ClaimsNamespace:      "http://tenant.info",
				Certificate:          []byte(certificateSuccessOrchestrateTest),
			},
			idToken,
			true,
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker, err := New(tt.cfg)
			require.NoError(t, err)

			ctx := authutils.WithAuthorization(context.Background(), fmt.Sprintf("Bearer %v", tt.token))
			checkedCtx, err := checker.Check(ctx)
			token := FromContext(checkedCtx)

			switch {
			case tt.isInvalid:
				assert.Error(t, err, "Check should have error")
				assert.Nil(t, token, "Token should not have been set")

			case !tt.isInvalid:
				assert.NoError(t, err, "Check should not have error")
				assert.NotNil(t, token, "Token should have been set")
			}
		})
	}
}
