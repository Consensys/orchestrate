package generator

import (
	"crypto/rsa"
	"fmt"
	"testing"
	"time"
)

var (
	rsaPrivateKeyLocation = "../../../../tests/certificates/e2e-tests.orchestrate.key.pem"
	rsaPrivateKeyString   = "MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCjQ2paojdNKLW44G8JSq5yVajgtyxhPbW01RdDxjKpVhFKyIG/5kpkEKMPxeSUs+KHjSQVJ76voM1xrn8sf+iRAh1zyexCD4EsKRd4fY0W/5Bi09zMDDhiFLoZXEmTe72x+nw1q9YFxBsPcwwOl6Ew6+FKXTOg/9mUfoSsoc22hCBSk81VAbKw1TbUqkKWYgMP8vV2uxur6sWIOtxVZo4vAbGvxIF3/RvieC1zAfXpElrFwhddPkBWJrUI13OZTn8AiDc5+DO8qfkAasxkGyYYCILKfj8/c+e4XzX9Ye6STkRMufM4f68LZ6VKZnzLzQcIZ6o6JDWuzrov3frBfxwPAgMBAAECggEARNLHg7t8SoeNy4i45hbYYRRhI5G0IK3t6nQl4YkslBvXIEpT//xpgbNNufl3OYR3SyMhgdWGWe0Ujga8T5sABBj7J3OIp/R3RJFx9nYewwIq8K5VFqNUJWyNYuF3lreEKQHp2Io+p6GasrGR9JjQ95mIGFwfxo/0Pdfzv/5ZhMWTmSTcOi504Vger5TaPobPFOnULq4y1A4eX4puiHDtvx09DUAWbAjGHpCYZjDGRdSXQArYQmUOKy7R46qKT/ollGOWivnEOgsFmXuUWs/shmcrDG4cGBkRrkxyIZhpnpNEEF5TYgulMMzwM+314e8W0lj9iiSB2nXzt8JhEwTz8QKBgQDSCouFj2lNSJDg+kz70eWBF9SQLrBTZ8JcMte3Q+CjCL1FpSVYYBRzwJNvWFyNNv7kHhYefqfcxUVSUnQ1eZIqTXtm9BsLXnTY+uEkV92spjVmfzBKZvtN3zzip97sfMT9qeyagHEHwpP+KaR0nyffAK+VPhlwNMKgQ9rzP4je+QKBgQDG/JwVaL2b53vi9CNh2XI8KNUd6rx6NGC6YTZ/xKVIgczGKTVex/w1DRWFTb0neUsdus5ITqaxQJtJDw/pOwoIag7Q0ttlLNpYsurx3mgMxpYY12/wurvp1NoU3Dq6ob7igfowP+ahUBchRwt1tlezn3TYxVoZpu9dZHtoynOtRwKBgB9vFJJYdBns0kHZM8w8DWzUdCtf0WOqE5xYv4/dyLCdjjXuETi4qFbqayYuwysfH+Zj2kuWCOkxXL6FOH8IQqeyENXHkoSRDkuqwCcAP1ynQzajskZwQwvUbPg+x039Hj4YQCCfOEtBA4T2Fnadmwn0wFJFiOkR/E6f2RSuXX2BAoGALvVqODsxk9s7B0IqH2tbZAsW0CqXNBesRA+w9tIHV2caViFfcPCs+jAORhkkbG5ZZbix+apl+CqQ+trNHHNMWNP+jxVTpTrChHAktdOQpoMu5MnipuLKedI7bPTT/zsweu/FhSFvYd4utzG26J6Rb9hPkOBx9N/KWTXfUcmFJv0CgYAUYVUvPe7MHSd5m8MulxRnVirWzUIUL9Pf1RKWOUq7Ue4oMxzE8CZCJstunCPWgyyxYXgj480PdIuL92eTR+LyaUESb6szZQTxaJfu0mEJS0KYWlONz+jKM4oC06dgJcCMvhgjta2KpXCm3qL1pmKwfFbOLWYBe5uMoHIn9FdJFQ=="
)

func TestGenerateAccessToken(t *testing.T) {
	type fields struct {
		claimsNamespace string
		privateKey      *rsa.PrivateKey
	}

	type args struct {
		customClaims map[string]interface{}
	}

	key, _ := LoadRSAPrivateKeyFromFile(rsaPrivateKeyLocation)

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantToken string
		wantErr   bool
	}{
		{
			"nominal case one line customClaim",
			fields{
				"http://tenant.info/",
				key,
			},
			args{
				customClaims: map[string]interface{}{"http://tenant.info/tenant_id": "f30c452b-e5fb-4102-a45d-bc00a060bcc6"},
			},
			"",
			false,
		},
		{
			"nominal case struct customClaim",
			fields{
				"http://tenant.info/",
				key,
			},
			args{
				customClaims: map[string]interface{}{"http://tenant.info/": map[string]interface{}{"tenant_id": "f30c452b-e5fb-4102-a45d-bc00a060bcc6"}},
			},
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JWTGenerator{
				tt.fields.claimsNamespace,
				tt.fields.privateKey,
			}
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

	type fields struct {
		multitenancy    bool
		claimsNamespace string
		privateKey      *rsa.PrivateKey
	}

	key, _ := LoadRSAPrivateKeyFromFile(rsaPrivateKeyLocation)

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantToken string
		wantErr   bool
	}{
		{
			"nominal case one line customClaim",
			fields{
				true,
				"http://tenant.info/",
				key,
			},
			args{
				customClaims: map[string]interface{}{"http://tenant.info/tenant_id": "f30c452b-e5fb-4102-a45d-bc00a060bcc6"},
			},
			"",
			false,
		},
		{
			"nominal case struct customClaim",
			fields{
				true,
				"http://tenant.info/",
				key,
			},
			args{
				customClaims: map[string]interface{}{"http://tenant.info/": map[string]interface{}{"tenant_id": "f30c452b-e5fb-4102-a45d-bc00a060bcc6"}},
			},
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JWTGenerator{
				ClaimsNamespace: tt.fields.claimsNamespace,
				privateKey:      tt.fields.privateKey,
			}
			gotToken, err := j.GenerateIDToken(tt.args.customClaims)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateIDToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("ID Token value: \n%s\n", gotToken)

		})
	}
}

func TestLoadRsaPrivateKeyFromVar(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			"nominal case",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, _ := LoadRsaPrivateKeyFromVar(rsaPrivateKeyString)
			if (got == nil) != tt.wantErr {
				t.Errorf("LoadRsaPrivateKeyFromVar() got is empty")
				return
			}
		})
	}
}
