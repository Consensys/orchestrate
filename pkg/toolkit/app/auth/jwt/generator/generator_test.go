package generator

import (
	"fmt"
	"testing"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/tls/certificate"
	tlstestutils "github.com/ConsenSys/orchestrate/pkg/toolkit/tls/testutils"
	"github.com/stretchr/testify/require"
)

func TestGenerateAccessToken(t *testing.T) {

	type args struct {
		customClaims map[string]interface{}
	}

	tests := []struct {
		name      string
		cfg       *Config
		args      args
		wantToken string
		wantErr   bool
	}{
		{
			"nominal case one line customClaim",
			&Config{
				ClaimsNamespace: "http://orchestrate.info/",
				KeyPair: &certificate.KeyPair{
					Key: []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
			},
			args{
				customClaims: map[string]interface{}{"http://orchestrate.info/tenant_id": "f30c452b-e5fb-4102-a45d-bc00a060bcc6"},
			},
			"",
			false,
		},
		{
			"nominal case struct customClaim",
			&Config{
				ClaimsNamespace: "http://orchestrate.info/",
				KeyPair: &certificate.KeyPair{
					Key: []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
			},
			args{
				customClaims: map[string]interface{}{"http://orchestrate.info": map[string]interface{}{"tenant_id": "f30c452b-e5fb-4102-a45d-bc00a060bcc6"}},
			},
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, err := New(tt.cfg)
			require.NoError(t, err)

			gotToken, err := j.GenerateAccessToken(tt.args.customClaims, 24*time.Hour)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("Access Token value: \n%s\n", gotToken)
		})
	}
}

func TestGenerateIDToken(t *testing.T) {
	type args struct {
		customClaims map[string]interface{}
	}

	tests := []struct {
		name      string
		cfg       *Config
		args      args
		wantToken string
		wantErr   bool
	}{
		{
			"nominal case one line customClaim",
			&Config{
				ClaimsNamespace: "http://orchestrate.info/",
				KeyPair: &certificate.KeyPair{
					Key: []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
			},
			args{
				customClaims: map[string]interface{}{"http://orchestrate.info/tenant_id": "f30c452b-e5fb-4102-a45d-bc00a060bcc6"},
			},
			"",
			false,
		},
		{
			"nominal case struct customClaim",
			&Config{
				ClaimsNamespace: "http://orchestrate.info/",
				KeyPair: &certificate.KeyPair{
					Key: []byte(tlstestutils.OneLineRSAKeyPEMA),
				},
			},
			args{
				customClaims: map[string]interface{}{"http://orchestrate.info": map[string]interface{}{"tenant_id": "f30c452b-e5fb-4102-a45d-bc00a060bcc6"}},
			},
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, err := New(tt.cfg)
			require.NoError(t, err)

			gotToken, err := j.GenerateIDToken(tt.args.customClaims)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateIDToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("UUID Token value: \n%s\n", gotToken)

		})
	}
}
