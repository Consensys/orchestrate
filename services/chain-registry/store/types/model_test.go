package types

import (
	"reflect"
	"testing"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/tls"
)

func TestBuildConfiguration(t *testing.T) {

	testSet := []struct {
		configs        []*Config
		expectedOutput func(*dynamic.Configuration) *dynamic.Configuration
		expectedError  bool
	}{
		{
			[]*Config{
				{Name: "testRoute", TenantID: "defaultTenantId", ConfigType: HTTPROUTER, Config: []byte(`{"service":"testService"}`)},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				c.HTTP.Routers["defaultTenantId-testRoute"] = &dynamic.Router{Service: "testService"}
				return c
			},
			false,
		},
		{
			[]*Config{
				{Name: "testMiddleware", TenantID: "defaultTenantId", ConfigType: HTTPMIDDLEWARE, Config: []byte(`{"addPrefix":{"prefix":"testPrefix"}}`)},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				c.HTTP.Middlewares["defaultTenantId-testMiddleware"] = &dynamic.Middleware{AddPrefix: &dynamic.AddPrefix{Prefix: "testPrefix"}}
				return c
			},
			false,
		},
		{
			[]*Config{
				{Name: "testServices", TenantID: "defaultTenantId", ConfigType: HTTPSERVICE, Config: []byte(`{"loadBalancer":{"servers":[{"url":"testUrl"}]}}`)},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				c.HTTP.Services["defaultTenantId-testServices"] = &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{{URL: "testUrl"}}}}
				return c
			},
			false,
		},
		{
			[]*Config{
				{Name: "testTCPRouter", TenantID: "defaultTenantId", ConfigType: TCPROUTER, Config: []byte(`{"service":"testService"}`)},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				c.TCP.Routers["defaultTenantId-testTCPRouter"] = &dynamic.TCPRouter{Service: "testService"}
				return c
			},
			false,
		},
		{
			[]*Config{
				{Name: "testTCPService", TenantID: "defaultTenantId", ConfigType: TCPSERVICE, Config: []byte(`{"loadBalancer":{"terminationDelay":10}}`)},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				terminationDelay := 10
				c.TCP.Services["defaultTenantId-testTCPService"] = &dynamic.TCPService{LoadBalancer: &dynamic.TCPServersLoadBalancer{TerminationDelay: &terminationDelay}}
				return c
			},
			false,
		},
		{
			[]*Config{
				{ConfigType: TLSCERTIFICATE, Config: []byte(`{"certFile":"testCertFile","keyFile":"testKeyFile"}`)},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				certificate := &tls.CertAndStores{
					Certificate: tls.Certificate{CertFile: "testCertFile", KeyFile: "testKeyFile"},
				}
				c.TLS.Certificates = append(c.TLS.Certificates, certificate)
				return c
			},
			true,
		},
		{
			[]*Config{
				{Name: "testTLSOptions", TenantID: "defaultTenantId", ConfigType: TLSOPTIONS, Config: []byte(`{"minVersion":"testMinVersion"}`)},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				c.TLS.Options["defaultTenantId-testTLSOptions"] = tls.Options{MinVersion: "testMinVersion"}
				return c
			},
			false,
		},
		{
			[]*Config{
				{Name: "testTLSStore", TenantID: "defaultTenantId", ConfigType: TLSSTORES, Config: []byte(`{"defaultCertificate":{"certFile":"testCertFile","keyFile":"testKeyFile"}}`)},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				c.TLS.Stores["defaultTenantId-testTLSStore"] = tls.Store{DefaultCertificate: &tls.Certificate{CertFile: "testCertFile", KeyFile: "testKeyFile"}}
				return c
			},
			false,
		},
		{
			[]*Config{
				{Name: "testRoute", TenantID: "defaultTenantId", ConfigType: HTTPROUTER, Config: []byte(`{"service":"testService"}`)},
				{Name: "testMiddleware", TenantID: "defaultTenantId", ConfigType: HTTPMIDDLEWARE, Config: []byte(`{"addPrefix":{"prefix":"testPrefix"}}`)},
				{Name: "testServices", TenantID: "defaultTenantId", ConfigType: HTTPSERVICE, Config: []byte(`{"loadBalancer":{"servers":[{"url":"testUrl"}]}}`)},
				{Name: "testTCPRouter", TenantID: "defaultTenantId", ConfigType: TCPROUTER, Config: []byte(`{"service":"testService"}`)},
				{Name: "testTCPService", TenantID: "defaultTenantId", ConfigType: TCPSERVICE, Config: []byte(`{"loadBalancer":{"terminationDelay":10}}`)},
				{Name: "testTLSOptions", TenantID: "defaultTenantId", ConfigType: TLSOPTIONS, Config: []byte(`{"minVersion":"testMinVersion"}`)},
				{Name: "testTLSStore", TenantID: "defaultTenantId", ConfigType: TLSSTORES, Config: []byte(`{"defaultCertificate":{"certFile":"testCertFile","keyFile":"testKeyFile"}}`)},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				c.HTTP.Routers["defaultTenantId-testRoute"] = &dynamic.Router{Service: "testService"}
				c.HTTP.Middlewares["defaultTenantId-testMiddleware"] = &dynamic.Middleware{AddPrefix: &dynamic.AddPrefix{Prefix: "testPrefix"}}
				c.HTTP.Services["defaultTenantId-testServices"] = &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{{URL: "testUrl"}}}}
				c.TCP.Routers["defaultTenantId-testTCPRouter"] = &dynamic.TCPRouter{Service: "testService"}
				terminationDelay := 10
				c.TCP.Services["defaultTenantId-testTCPService"] = &dynamic.TCPService{LoadBalancer: &dynamic.TCPServersLoadBalancer{TerminationDelay: &terminationDelay}}
				c.TLS.Options["defaultTenantId-testTLSOptions"] = tls.Options{MinVersion: "testMinVersion"}
				c.TLS.Stores["defaultTenantId-testTLSStore"] = tls.Store{DefaultCertificate: &tls.Certificate{CertFile: "testCertFile", KeyFile: "testKeyFile"}}
				return c
			},
			false,
		},
		{
			[]*Config{
				{Name: "testTLSStore", ConfigType: UNKNOWN, Config: []byte(`{"defaultCertificate":{"certFile":"testCertFile","keyFile":"testKeyFile"}}`)},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				return c
			},
			true,
		},
		{
			[]*Config{
				{Name: "testTLSStore", ConfigType: HTTPROUTER, Config: []byte(`{"wrongKey":"testRoute"}`)},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				return c
			},
			true,
		},
	}

	for i, test := range testSet {
		output, err := BuildConfiguration(test.configs)
		if (err == nil && test.expectedError) || (err != nil && !test.expectedError) {
			t.Errorf("Chain-registry - Store: Expecting the following error %v but got %v", test.expectedError, err)
		} else if err != nil && test.expectedError {
			continue
		}

		expectedOutput := test.expectedOutput(NewTraefikConfig())
		eq := reflect.DeepEqual(expectedOutput, output)
		if !eq {
			t.Errorf("Chain-registry - Store (%d/%d): expected %v but got %v", i+1, len(testSet), expectedOutput, output)
		}

	}
}

func TestTraefikConfigToStoreConfig(t *testing.T) {

	testSet := []struct {
		configs        *dynamic.Configuration
		tenantID       string
		expectedOutput []*Config
		expectedError  bool
	}{
		{
			&dynamic.Configuration{
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"testRoute": {Service: "testService"},
					},
				},
			},
			"testTenant",
			[]*Config{
				{Name: "testRoute", TenantID: "testTenant", ConfigType: HTTPROUTER, Config: []byte(`{"service":"testService"}`)},
			},
			false,
		},
		{
			&dynamic.Configuration{
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"testRoute":  {Service: "testService"},
						"testRoute2": {Service: "testService2"},
					},
				},
			},
			"testTenant",
			[]*Config{
				{Name: "testRoute", TenantID: "testTenant", ConfigType: HTTPROUTER, Config: []byte(`{"service":"testService"}`)},
				{Name: "testRoute2", TenantID: "testTenant", ConfigType: HTTPROUTER, Config: []byte(`{"service":"testService2"}`)},
			},
			false,
		},
		{
			&dynamic.Configuration{
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"testRoute": {Service: "testService"},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"testTCPRouter": {EntryPoints: []string{"testEntryPoints"}},
					},
				},
			},
			"testTenant",
			[]*Config{
				{Name: "testRoute", TenantID: "testTenant", ConfigType: HTTPROUTER, Config: []byte(`{"service":"testService"}`)},
				{Name: "testTCPRouter", TenantID: "testTenant", ConfigType: TCPROUTER, Config: []byte(`{"entryPoints":["testEntryPoints"]}`)},
			},
			false,
		},
		{
			&dynamic.Configuration{
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"testRoute": {Service: "testService"},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"testTCPRouter": {EntryPoints: []string{"testEntryPoints"}},
					},
				},
				TLS: &dynamic.TLSConfiguration{
					Certificates: []*tls.CertAndStores{
						{Certificate: tls.Certificate{CertFile: "testCertFile", KeyFile: "testKeyFile"}},
						{Certificate: tls.Certificate{CertFile: "testCertFile2", KeyFile: "testKeyFile2"}},
					},
				},
			},
			"testTenant",
			[]*Config{
				{Name: "testRoute", TenantID: "testTenant", ConfigType: HTTPROUTER, Config: []byte(`{"service":"testService"}`)},
				{Name: "testTCPRouter", TenantID: "testTenant", ConfigType: TCPROUTER, Config: []byte(`{"entryPoints":["testEntryPoints"]}`)},
				{TenantID: "testTenant", ConfigType: TLSCERTIFICATE, Config: []byte(`{"certFile":"testCertFile","keyFile":"testKeyFile"}`)},
				{TenantID: "testTenant", ConfigType: TLSCERTIFICATE, Config: []byte(`{"certFile":"testCertFile2","keyFile":"testKeyFile2"}`)},
			},
			true,
		},
		{
			&dynamic.Configuration{
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"testRoute": {Service: "testService"},
					},
				},
			},
			"testTenant",
			[]*Config{
				{Name: "testRoute", TenantID: "testTenant", ConfigType: HTTPROUTER, Config: []byte(`{"service":"testService"}`)},
			},
			false,
		},
	}

	for i, test := range testSet {
		output, err := TraefikConfigToStoreConfig(test.configs, test.tenantID)
		if (err == nil && test.expectedError) || (err != nil && !test.expectedError) {
			t.Errorf("Chain-registry - Store: Expecting the following error %v but got %v", test.expectedError, err)
		} else if err != nil && test.expectedError {
			continue
		}

		eq := reflect.DeepEqual(test.expectedOutput, output)
		if !eq {
			t.Errorf("Chain-registry - Store (%d/%d): expected %v but got %v", i+1, len(testSet), test.expectedOutput, output)
		}
	}
}
