// +build !race

package multitenancy

import (
	"fmt"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine/testutils"
	errors "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/token"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

const (
	idToken                           = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6Ik1UTXpRamRETUROR01qazVOREpCTUVKQk56ZzBSVUk0UWpjek5UZzJNMFUxT0VZNVJrUTNRdyJ9.eyJodHRwOi8vdGVuYW50LmluZm8vdGVuYW50X2lkIjoiYjQ5ZWUxYmMtZjBmYS00MzBkLTg5YjItYTRmZDBkYzk4OTA2IiwiaHR0cDovL3RlbmFudC5pbmZvL3RlbmFudF9yb2xlIjoidXNlciIsImh0dHA6Ly90ZW5hbnQuaW5mby90ZW5hbnRfY29tcGFnbnkiOiJQZWdhU3lzIiwibmlja25hbWUiOiJmb28iLCJuYW1lIjoiZm9vQGJhci5jb20iLCJwaWN0dXJlIjoiaHR0cHM6Ly9zLmdyYXZhdGFyLmNvbS9hdmF0YXIvZjNhZGE0MDVjZTg5MGI2ZjgyMDQwOTRkZWIxMmQ4YTg_cz00ODAmcj1wZyZkPWh0dHBzJTNBJTJGJTJGY2RuLmF1dGgwLmNvbSUyRmF2YXRhcnMlMkZmby5wbmciLCJ1cGRhdGVkX2F0IjoiMjAxOS0xMi0wNlQwOTo0ODowMS41NTNaIiwiaXNzIjoiaHR0cHM6Ly9kZXYtYmQ2ZTNqbGMuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfDVkZGU4ZTYyNzY5YTJkMGVkM2FmNTM4ZSIsImF1ZCI6IlpDZTdKdUNsaXUyMFIwc0xwU0UwdzhJN1d3YTE2MldkIiwiaWF0IjoxNTc1NjI1NjgyLCJleHAiOjE1NzU3MTIwODJ9.muHMxGe0EaSYnRCVpVAPeIfeEr4VLnN54DcWOxk6CMBUlNq2gzElxiKkZ2IUS6oZXCwHvob40mMJQJyIPpRBn23ZsIZLK3Iy4Xbf-TytvtSKWMX4Jiw1WgNey7_DsjHtT6Wi9OufS2NF49sK39m0hDXf2GCqqtYFg5XNQLMujfDdplxN2gRHP3VEey3PtSMBFIdlAkv2mCA5SPBlxmkCtGmgiQa223bPl2rnCA5PF7XjNVTg2v59m34ADZ8cR-J6h1UrKPXFmCXEO1gHC_wpiN7E0pjjnJVORDN27b5zAASADPSh9tyZlWbZa14SAP8M9gzOChS5z5b31efuvA8Rxw"
	accessTokenUntrustedSigner        = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IlJFRTVSRUpETURFeFJrVTFNVGhDUTBGQ05rTXdSVEkyTVVSRk1qQXlOekUyTjBNMU1rWXpNZyJ9.eyJodHRwczovL2Rldi5hYmNjZC5jb20iOnsicm9sZXMiOlsiVHJhZGUgRXhlY3V0b3IiXSwibGVnYWxFbnRpdGllcyI6WzFdLCJsZWdhbEVudGl0eSI6MX0sImlzcyI6Imh0dHBzOi8vZGV2LWFiY2NkLmV1LmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1ZDhjODFmMTRiMmZlYjBkNzA0YzVlMGYiLCJhdWQiOlsiYmVhdC1hdXRoMC1hcGkiLCJodHRwczovL2Rldi1hYmNjZC5ldS5hdXRoMC5jb20vdXNlcmluZm8iXSwiaWF0IjoxNTc1NTQxMjA1LCJleHAiOjE1NzU2Mjc2MDUsImF6cCI6ImVkQ3hvVlRPcFhCN2t0NmdpVFdoOG1CZ0pnTVdvTzJ2Iiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInBlcm1pc3Npb25zIjpbImNhOmNyZWF0ZSIsImNhOnVwZGF0ZSIsImNhOnZpZXdEZXRhaWxzIiwibGVnYWwtZW50aXR5OnJlYWQiLCJtYXN0ZXItZGF0YTpyZWFkIiwibm9taW5hdGlvbjpjcmVhdGUiLCJub21pbmF0aW9uOnVwZGF0ZSIsIm5vbWluYXRpb246dmlld0RldGFpbHMiLCJ2b3lhZ2U6Y3JlYXRlIiwidm95YWdlOnJlYWQiXX0.jJnJjTLHsElFU3O7xKuh7jL1ho9-Z7Jxco16hDxoRg_TFdOCN82wVeJHZbDjdLqjV0k4F05YWEmFWn7CEAmr43ndoprsAr3OfBnrjYKyJ4oqiguPAUakBqoLtaEE-AsxyQmCzZGwKXHtNMDIhh0vwHVASdHTwxiApumRWfEXzmu5pmOYwoTJ8vVSUVCCDG3hL6u4UxYdng30XlWgbn_Szlaq9sIoIllZOL8vn4hkkW98CQfjexpaYDjywVfbPD3-TSSznHiF6TvmogCttkb73hbJF246hq-guR0nfdQm1ivAUkzXcUOql6QtHvYgdrzw5xPOqNIMihFvIK8XRCZ_pw"
	accessTokenWithoutTenantID        = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6Ik1UTXpRamRETUROR01qazVOREpCTUVKQk56ZzBSVUk0UWpjek5UZzJNMFUxT0VZNVJrUTNRdyJ9.eyJpc3MiOiJodHRwczovL2Rldi1iZDZlM2psYy5hdXRoMC5jb20vIiwic3ViIjoiY1hFTDNyS0NIdnhPbTV2QmE5TU1hblVpRVNuaUxjOE1AY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vZGV2LWJkNmUzamxjLmF1dGgwLmNvbS9hcGkvdjIvIiwiaWF0IjoxNTc0ODY4Nzc3LCJleHAiOjE1NzQ5NTUxNzcsImF6cCI6ImNYRUwzcktDSHZ4T201dkJhOU1NYW5VaUVTbmlMYzhNIiwiZ3R5IjoiY2xpZW50LWNyZWRlbnRpYWxzIn0.PMgDjxiL4Hddaxrsj8FyEbyN-aeXZBAOhaIZR0tsnHYVhe6Xbdm-MBvAreKeiqY_WJmJdw9GZxNOBeT9Tk9WPjojArN4gIllyym4OnRrBAdDx0KR-Lg4gNAXUMWazEP1FQbBRhXWbMFASxlyr8I6Evzel55MBLmClTpD6kR_Z2QY8JJRAw21i55GjeWTAa-9NtMWzWl7klzNjSyNGGwD3hpcYbzUUjAdU4IJ7LPZ2MVceFbjEUuf1vz8PTE54W8caxgYXoiismxArG5Ck_KCYoLHT-PtgSeXWwxYCimejt-QwgquYtFzjOybUSanGu1BCVzBAUGiblNLDEmU-eD_Rg"
	certificateOneLineOrchestrateTest = "MIIDBzCCAe+gAwIBAgIJCOOsj4KofbjsMA0GCSqGSIb3DQEBCwUAMCExHzAdBgNVBAMTFmRldi1iZDZlM2psYy5hdXRoMC5jb20wHhcNMTkxMTI2MTYzODMwWhcNMzMwODA0MTYzODMwWjAhMR8wHQYDVQQDExZkZXYtYmQ2ZTNqbGMuYXV0aDAuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApWBAkbQrPMOeF7GFz9EhKbsUFOg3WxVtPlvMtjkTtgxJe5ke5dxc2F9YeMB+1N+I2ozQa1ReCWAun4rGz4ovjxI4PeUT0exFbI4oKd2bKOEd/IVGmabgUEm3FlSSq0jOEgu8JMpmGZIEGi3RMg8E1mAIJf5VwiIrCE6sP7IY9wrBaavmMdJ/i2a0gmjmPNqD8Y2bMi0fWW5frmGibMPEaddG8/Daj3SMWo8N8nhW1VX3JyQcuA3Jxvsyj8aYudoCWIhbYSsdeVY3JmUnIcGZ7XVJH7COEwPmnxQ5uJAnqPfbItPMN9yzGqYxC4eC3UGzKJE5dfOcLCDJOe6AtKuxmwIDAQABo0IwQDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBQxhyU5rj46P2H8VwI5Rq/nwsHSNTAOBgNVHQ8BAf8EBAMCAoQwDQYJKoZIhvcNAQELBQADggEBAG4WbRfOYeUNz637G5eFC3LMGa3bu+S/ln+NON3ZI49adCxcXElR8fIpXdtq/HzyZGcWfdo5+sgaSKRAD4iWdEFtPkK840gIdFXf7lScSBo76uqiMvbw1xGbyNcsNbUppTM1FmfrJ25CaMGG+9yd8gjBuHNLOmZXGkvo9et0ECKQEku9BunuGwIdWTaq5BTEufqby4sEtv0ZwLgSwsooMRCMUIU2e/MM9wyD21Gc9Qp2v3/TI2282eVrIWunWE0WgMG0KlIdfFuGpGqJUfXjBVD+WAvV/E2lFraILs7sIp8U35hmJq4vG0kjG9B+JKHYswyLtnw+3LVuAbUNiB5MLM4="
)

var testsNumError = 3
var testsNum = 2
var tenantPath = "http://tenant.info/"

var keyExpectedValue = "expected-test"

type MultiTenancyTestSuite struct {
	testutils.HandlerTestSuite
	authManager *token.AuthToken
}

func makeMultiTenancyContextForError(i int) *engine.TxContext {
	ctx := engine.NewTxContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())
	ctx.Envelope.Receipt = &ethereum.Receipt{}

	_ = os.Setenv("AUTH_SERVICE_CERTIFICATE", certificateOneLineOrchestrateTest)

	switch i % 3 {
	case 0:
		// Error Use case:  Token is expired
		ctx.Envelope.SetMetadataValue(token.OauthToken, idToken)
		ctx.Set("errors", 1)
		ctx.Set("error.code", errors.Unauthenticated)
	case 1:
		// Error Use case:  UntrustedSigner
		ctx.Envelope.SetMetadataValue(token.OauthToken, accessTokenUntrustedSigner)
		ctx.Set("errors", 1)
		ctx.Set("error.code", errors.Unauthenticated)
	case 2:
		// Error Use case:  Token no present
		ctx.Set("errors", 1)
		ctx.Set("error.code", errors.NotFound)
	default:
		panic(fmt.Sprintf("No test case with number %d", i))
	}
	return ctx
}

func (m *MultiTenancyTestSuite) TestMultiTenancyErrorUseCases() {
	m.authManager = token.New()
	m.Handler = ExtractTenant(m.authManager)

	var txctxs []*engine.TxContext
	for i := 0; i < testsNumError; i++ {
		txctxs = append(txctxs, makeMultiTenancyContextForError(i))
	}

	// Handle contexts
	m.Handle(txctxs)

	for _, txctx := range txctxs {
		assert.Len(m.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "Expected right count of errors", txctx.Envelope.Args)
		for _, err := range txctx.Envelope.Errors {
			assert.Equal(m.T(), txctx.Get("error.code").(uint64), err.GetCode(), "Error code be correct")
		}
	}
}

func makeMultiTenancyContext(i int) *engine.TxContext {
	ctx := engine.NewTxContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())
	ctx.Envelope.Receipt = &ethereum.Receipt{}

	_ = os.Setenv("AUTH_SERVICE_CERTIFICATE", certificateOneLineOrchestrateTest)
	_ = os.Setenv("TENANT_NAMESPACE", tenantPath)

	switch i % 2 {
	case 0:
		// Error Use case:  Token is expired
		ctx.Envelope.SetMetadataValue(token.OauthToken, idToken)
		ctx.Set(keyExpectedValue, "b49ee1bc-f0fa-430d-89b2-a4fd0dc98906")
		ctx.Set("errors", 0)
	case 1:
		// Error Use case:  UntrustedSigner
		ctx.Envelope.SetMetadataValue(token.OauthToken, accessTokenWithoutTenantID)
		ctx.Set("errors", 1)
		ctx.Set("error.code", errors.NotFound)
	default:
		panic(fmt.Sprintf("No test case with number %d", i))
	}
	return ctx
}

func (m *MultiTenancyTestSuite) TestMultiTenancy() {
	m.authManager = &token.AuthToken{
		Parser: &jwt.Parser{
			SkipClaimsValidation: true,
		},
	}
	m.Handler = ExtractTenant(m.authManager)

	var txctxs []*engine.TxContext
	for i := 0; i < testsNum; i++ {
		txctxs = append(txctxs, makeMultiTenancyContext(i))
	}

	// Handle contexts
	m.Handle(txctxs)

	for _, txctx := range txctxs {
		assert.Len(m.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "Expected right count of errors", txctx.Envelope.Args)
		if txctx.Get("errors").(int) != 0 {
			for _, err := range txctx.Envelope.Errors {
				assert.Equal(m.T(), txctx.Get("error.code").(uint64), err.GetCode(), "Error code be correct")
			}
		} else {
			assert.Equal(m.T(), txctx.Get(keyExpectedValue).(string), txctx.Get(TenantIDKey).(string), "Expected correct TenantIDKey")
		}
	}
}

func TestMultiTenancy(t *testing.T) {
	suite.Run(t, new(MultiTenancyTestSuite))
}
