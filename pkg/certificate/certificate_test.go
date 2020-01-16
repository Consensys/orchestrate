package certificate

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
)

const (
	certificateBadFormatedOrchestrateTest = "MIIDBzCCAe+gAwIBAgIJCOOsj4KofbjsMA0GCSqGSIb3DQEBCwUAMCExHzAdBgNVBAMTFmRldi1iZDZlM2psYy5hdXRoMC5jb20wHhcNMTkxMTI2MTYzODMwWhcNMzMwODA0MTYzODMwWjAhMR8wHQYDVQQDExZkZXYtYmQ2ZTNqbGMuYXV0aDAuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApWBAkbQrPMOeF7GFz9EhKbsUFOg3WxVtPlvMtjkTtgxJe5ke5dxc2F9YeMB+1N+I2ozQa1ReCWAun4rGz4ovjxI4PeUT0exFbI4oKd2bKOEd/IVGmabgUEm3FlSSq0jOEgu8JMpmGZIEGi3RMg8E1mAIJf5VwiIrCE6sP7IY9wrBaavmMdJ/i2a0gmjmPNqD8Y2bMi0fWW5frmGibMPEaddG8Daj3SMWo8N8nhW1VX3JyQcuA3Jxvsyj8aYudoCWIhbYSsdeVY3JmUnIcGZ7XVJH7COEwPmnxQ5uJAnqPfbItPMN9yzGqYxC4eC3UGzKJE5dfOcLCDJOe6AtKuxmwIDAQABo0IwQDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBQxhyU5rj46P2H8VwI5Rq/nwsHSNTAOBgNVHQ8BAf8EBAMCAoQwDQYJKoZIhvcNAQELBQADggEBAG4WbRfOYeUNz637G5eFC3LMGa3bu+S/ln+NON3ZI49adCxcXElR8fIpXdtq/HzyZGcWfdo5+sgaSKRAD4iWdEFtPkK840gIdFXf7lScSBo76uqiMvbw1xGbyNcsNbUppTM1FmfrJ25CaMGG+9yd8gjBuHNLOmZXGkvo9et0ECKQEku9BunuGwIdWTaq5BTEufqby4sEtv0ZwLgSwsooMRCMUIU2e/MM9wyD21Gc9Qp2v3/TI2282eVrIWunWE0WgMG0KlIdfFuGpGqJUfXjBVD+WAvV/E2lFraILs7sIp8U35hmJq4vG0kjG9B+JKHYswyLtnw+3LVuAbUNiB5MLM4="
	certificateOneLineOrchestrateTest     = "MIIDBzCCAe+gAwIBAgIJCOOsj4KofbjsMA0GCSqGSIb3DQEBCwUAMCExHzAdBgNVBAMTFmRldi1iZDZlM2psYy5hdXRoMC5jb20wHhcNMTkxMTI2MTYzODMwWhcNMzMwODA0MTYzODMwWjAhMR8wHQYDVQQDExZkZXYtYmQ2ZTNqbGMuYXV0aDAuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApWBAkbQrPMOeF7GFz9EhKbsUFOg3WxVtPlvMtjkTtgxJe5ke5dxc2F9YeMB+1N+I2ozQa1ReCWAun4rGz4ovjxI4PeUT0exFbI4oKd2bKOEd/IVGmabgUEm3FlSSq0jOEgu8JMpmGZIEGi3RMg8E1mAIJf5VwiIrCE6sP7IY9wrBaavmMdJ/i2a0gmjmPNqD8Y2bMi0fWW5frmGibMPEaddG8/Daj3SMWo8N8nhW1VX3JyQcuA3Jxvsyj8aYudoCWIhbYSsdeVY3JmUnIcGZ7XVJH7COEwPmnxQ5uJAnqPfbItPMN9yzGqYxC4eC3UGzKJE5dfOcLCDJOe6AtKuxmwIDAQABo0IwQDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBQxhyU5rj46P2H8VwI5Rq/nwsHSNTAOBgNVHQ8BAf8EBAMCAoQwDQYJKoZIhvcNAQELBQADggEBAG4WbRfOYeUNz637G5eFC3LMGa3bu+S/ln+NON3ZI49adCxcXElR8fIpXdtq/HzyZGcWfdo5+sgaSKRAD4iWdEFtPkK840gIdFXf7lScSBo76uqiMvbw1xGbyNcsNbUppTM1FmfrJ25CaMGG+9yd8gjBuHNLOmZXGkvo9et0ECKQEku9BunuGwIdWTaq5BTEufqby4sEtv0ZwLgSwsooMRCMUIU2e/MM9wyD21Gc9Qp2v3/TI2282eVrIWunWE0WgMG0KlIdfFuGpGqJUfXjBVD+WAvV/E2lFraILs7sIp8U35hmJq4vG0kjG9B+JKHYswyLtnw+3LVuAbUNiB5MLM4="
	certificateExpectedOrchestrateTest    = `
	MIIDBzCCAe+gAwIBAgIJCOOsj4KofbjsMA0GCSqGSIb3DQEBCwUAMCExHzAdBgNV
	BAMTFmRldi1iZDZlM2psYy5hdXRoMC5jb20wHhcNMTkxMTI2MTYzODMwWhcNMzMw 
	ODA0MTYzODMwWjAhMR8wHQYDVQQDExZkZXYtYmQ2ZTNqbGMuYXV0aDAuY29tMIIB 
	IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApWBAkbQrPMOeF7GFz9EhKbsU
	FOg3WxVtPlvMtjkTtgxJe5ke5dxc2F9YeMB+1N+I2ozQa1ReCWAun4rGz4ovjxI4
	PeUT0exFbI4oKd2bKOEd/IVGmabgUEm3FlSSq0jOEgu8JMpmGZIEGi3RMg8E1mAI
	Jf5VwiIrCE6sP7IY9wrBaavmMdJ/i2a0gmjmPNqD8Y2bMi0fWW5frmGibMPEaddG
	8/Daj3SMWo8N8nhW1VX3JyQcuA3Jxvsyj8aYudoCWIhbYSsdeVY3JmUnIcGZ7XVJ
	H7COEwPmnxQ5uJAnqPfbItPMN9yzGqYxC4eC3UGzKJE5dfOcLCDJOe6AtKuxmwID
	AQABo0IwQDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBQxhyU5rj46P2H8VwI5
	Rq/nwsHSNTAOBgNVHQ8BAf8EBAMCAoQwDQYJKoZIhvcNAQELBQADggEBAG4WbRfO
	YeUNz637G5eFC3LMGa3bu+S/ln+NON3ZI49adCxcXElR8fIpXdtq/HzyZGcWfdo5
	+sgaSKRAD4iWdEFtPkK840gIdFXf7lScSBo76uqiMvbw1xGbyNcsNbUppTM1Fmfr
	J25CaMGG+9yd8gjBuHNLOmZXGkvo9et0ECKQEku9BunuGwIdWTaq5BTEufqby4sE
	tv0ZwLgSwsooMRCMUIU2e/MM9wyD21Gc9Qp2v3/TI2282eVrIWunWE0WgMG0KlId
	fFuGpGqJUfXjBVD+WAvV/E2lFraILs7sIp8U35hmJq4vG0kjG9B+JKHYswyLtnw+
	3LVuAbUNiB5MLM4=`

	pubKeyExpectedOrchestrateTest = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApWBAkbQrPMOeF7GFz9Eh
KbsUFOg3WxVtPlvMtjkTtgxJe5ke5dxc2F9YeMB+1N+I2ozQa1ReCWAun4rGz4ov
jxI4PeUT0exFbI4oKd2bKOEd/IVGmabgUEm3FlSSq0jOEgu8JMpmGZIEGi3RMg8E
1mAIJf5VwiIrCE6sP7IY9wrBaavmMdJ/i2a0gmjmPNqD8Y2bMi0fWW5frmGibMPE
addG8/Daj3SMWo8N8nhW1VX3JyQcuA3Jxvsyj8aYudoCWIhbYSsdeVY3JmUnIcGZ
7XVJH7COEwPmnxQ5uJAnqPfbItPMN9yzGqYxC4eC3UGzKJE5dfOcLCDJOe6AtKux
mwIDAQAB
-----END RSA PUBLIC KEY-----
`
)

func TestDecodeStringToCertificate(t *testing.T) {
	tests := []struct {
		name    string
		param   string
		wantErr bool
	}{
		{
			"nominal case already formated",
			certificateExpectedOrchestrateTest,
			false,
		},
		{
			"nominal case need to be formated",
			certificateOneLineOrchestrateTest,
			false,
		},
		{
			"format error",
			certificateBadFormatedOrchestrateTest,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeStringToCertificate(tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeStringToCertificate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("DecodeStringToCertificate() got = %v, certificateBadFormatedOrchestrateTest different fron nil", got)
			}
		})
	}
}

func TestEncodePem(t *testing.T) {
	type args struct {
		pemString string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"nominal case",
			args{
				certificateBadFormatedOrchestrateTest,
			},
			pubKeyExpectedOrchestrateTest,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			block := &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: []byte(tt.args.pemString),
			}

			pemEncode := pem.EncodeToMemory(block)

			t.Logf("Pem encoded result: \n %s", string(pemEncode))
		})
	}
}

func TestDecodeAndFormat(t *testing.T) {
	type args struct {
		PemString string
		Type      string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"nominal case",
			args{
				certificateOneLineOrchestrateTest,
				certificateType,
			},
			pubKeyExpectedOrchestrateTest,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			block := &pem.Block{
				Type:  tt.args.Type,
				Bytes: []byte(tt.args.PemString),
			}

			var err error
			block.Bytes, err = decodeBase64(block.Bytes)
			if err != nil {
				t.Errorf("decodeBase64() got = %v\n", err)
			}

			pemEncode := pem.EncodeToMemory(block)
			t.Logf("Pem encoded result: \n %s", string(pemEncode))

			block, _ = pem.Decode(pemEncode)

			if block != nil {
				// Parse the certificate
				decodedCert, err := x509.ParseCertificate(block.Bytes)
				if err != nil {
					t.Errorf("ParseCertificate() got = %v\n", err)
				}

				t.Logf("Pem encoded result: \n %v", decodedCert)
			} else {
				t.Errorf("Can not decode PEM certificate:%v\n", string(pemEncode))
			}
		})
	}
}
