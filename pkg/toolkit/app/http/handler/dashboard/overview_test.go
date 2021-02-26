// +build unit

package dashboard

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/config/runtime"
	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/provider/docker"
	"github.com/containous/traefik/v2/pkg/provider/file"
	"github.com/containous/traefik/v2/pkg/provider/kubernetes/crd"
	"github.com/containous/traefik/v2/pkg/provider/kubernetes/ingress"
	"github.com/containous/traefik/v2/pkg/provider/marathon"
	"github.com/containous/traefik/v2/pkg/provider/rancher"
	"github.com/containous/traefik/v2/pkg/provider/rest"
	"github.com/containous/traefik/v2/pkg/tracing/jaeger"
	traefiktypes "github.com/containous/traefik/v2/pkg/types"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateExpected = flag.Bool("update_expected", false, "Update expected files in testdata")

func TestOverview(t *testing.T) {
	type expected struct {
		statusCode int
		jsonFile   string
	}

	testCases := []struct {
		desc      string
		path      string
		staticCfg *traefikstatic.Configuration
		infos     *runtime.Infos
		expected  expected
	}{
		{
			desc:      "without data in the dynamic configuration",
			path:      "/overview",
			staticCfg: &traefikstatic.Configuration{API: &traefikstatic.API{}, Global: &traefikstatic.Global{}},
			infos:     &runtime.Infos{},
			expected: expected{
				statusCode: http.StatusOK,
				jsonFile:   "testdata/overview-empty.json",
			},
		},
		{
			desc:      "with data in the dynamic configuration",
			path:      "/overview",
			staticCfg: &traefikstatic.Configuration{API: &traefikstatic.API{}, Global: &traefikstatic.Global{}},
			infos: &runtime.Infos{
				Services: map[string]*runtime.ServiceInfo{
					"foo-service@myprovider": {
						Service: &dynamic.Service{
							ReverseProxy: &dynamic.ReverseProxy{
								LoadBalancer: &dynamic.LoadBalancer{
									Servers: []*dynamic.Server{
										{
											URL: "http://127.0.0.1",
										},
									},
								},
							},
						},
						Status: runtime.StatusEnabled,
					},
					"bar-service@myprovider": {
						Service: &dynamic.Service{
							ReverseProxy: &dynamic.ReverseProxy{
								LoadBalancer: &dynamic.LoadBalancer{
									Servers: []*dynamic.Server{
										{
											URL: "http://127.0.0.1",
										},
									},
								},
							},
						},
						Status: runtime.StatusWarning,
					},
					"fii-service@myprovider": {
						Service: &dynamic.Service{
							ReverseProxy: &dynamic.ReverseProxy{
								LoadBalancer: &dynamic.LoadBalancer{
									Servers: []*dynamic.Server{
										{
											URL: "http://127.0.0.1",
										},
									},
								},
							},
						},
						Status: runtime.StatusDisabled,
					},
				},
				Middlewares: map[string]*runtime.MiddlewareInfo{
					"auth@myprovider": {
						Middleware: &dynamic.Middleware{
							Auth: &dynamic.Auth{},
						},
						Status: runtime.StatusEnabled,
					},
					"addPrefixTest@myprovider": {
						Middleware: &dynamic.Middleware{
							Middleware: &traefikdynamic.Middleware{
								AddPrefix: &traefikdynamic.AddPrefix{
									Prefix: "/titi",
								},
							},
						},
					},
					"addPrefixTest@anotherprovider": {
						Middleware: &dynamic.Middleware{
							Middleware: &traefikdynamic.Middleware{
								AddPrefix: &traefikdynamic.AddPrefix{
									Prefix: "/toto",
								},
							},
						},
						Err:    []string{"error"},
						Status: runtime.StatusDisabled,
					},
				},
				Routers: map[string]*runtime.RouterInfo{
					"bar@myprovider": {
						Router: &dynamic.Router{
							Router: &traefikdynamic.Router{
								EntryPoints: []string{"web"},
								Service:     "foo-service@myprovider",
								Rule:        "Host(`foo.bar`)",
								Middlewares: []string{"auth", "addPrefixTest@anotherprovider"},
							},
						},
						Status: runtime.StatusEnabled,
					},
					"test@myprovider": {
						Router: &dynamic.Router{
							Router: &traefikdynamic.Router{
								EntryPoints: []string{"web"},
								Service:     "foo-service@myprovider",
								Rule:        "Host(`foo.bar.other`)",
								Middlewares: []string{"addPrefixTest", "auth"},
							},
						},
						Status: runtime.StatusWarning,
					},
					"foo@myprovider": {
						Router: &dynamic.Router{
							Router: &traefikdynamic.Router{
								EntryPoints: []string{"web"},
								Service:     "foo-service@myprovider",
								Rule:        "Host(`foo.bar.other`)",
								Middlewares: []string{"addPrefixTest", "auth"},
							},
						},
						Status: runtime.StatusDisabled,
					},
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				jsonFile:   "testdata/overview-dynamic.json",
			},
		},
		{
			desc: "with providers",
			path: "/overview",
			staticCfg: &traefikstatic.Configuration{
				Global: &traefikstatic.Global{},
				API:    &traefikstatic.API{},
				Providers: &traefikstatic.Providers{
					Docker:            &docker.Provider{},
					File:              &file.Provider{},
					Marathon:          &marathon.Provider{},
					KubernetesIngress: &ingress.Provider{},
					KubernetesCRD:     &crd.Provider{},
					Rest:              &rest.Provider{},
					Rancher:           &rancher.Provider{},
				},
			},
			infos: &runtime.Infos{},
			expected: expected{
				statusCode: http.StatusOK,
				jsonFile:   "testdata/overview-providers.json",
			},
		},
		{
			desc: "with features",
			path: "/overview",
			staticCfg: &traefikstatic.Configuration{
				Global: &traefikstatic.Global{},
				API:    &traefikstatic.API{},
				Metrics: &traefiktypes.Metrics{
					Prometheus: &traefiktypes.Prometheus{},
				},
				Tracing: &traefikstatic.Tracing{
					Jaeger: &jaeger.Config{},
				},
			},
			infos: &runtime.Infos{},
			expected: expected{
				statusCode: http.StatusOK,
				jsonFile:   "testdata/overview-features.json",
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			router := mux.NewRouter()
			handler := NewOverview(test.staticCfg, test.infos)
			handler.Append(router)
			server := httptest.NewServer(router)

			resp, err := http.DefaultClient.Get(server.URL + test.path)
			require.NoError(t, err)

			require.Equal(t, test.expected.statusCode, resp.StatusCode)

			if test.expected.jsonFile == "" {
				return
			}

			assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
			contents, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			err = resp.Body.Close()
			require.NoError(t, err)

			if *updateExpected {
				var results interface{}
				err = json.Unmarshal(contents, &results)
				require.NoError(t, err)

				newJSON, e := json.MarshalIndent(results, "", "\t")
				require.NoError(t, e)

				err = ioutil.WriteFile(test.expected.jsonFile, newJSON, 0644)
				require.NoError(t, err)
			}

			data, err := ioutil.ReadFile(test.expected.jsonFile)
			require.NoError(t, err)
			assert.JSONEq(t, string(data), string(contents))
		})
	}
}
