// +build unit

package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/runtime"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func TestHTTP(t *testing.T) {
	type expected struct {
		statusCode int
		nextPage   string
		jsonFile   string
	}

	testCases := []struct {
		desc     string
		path     string
		infos    *runtime.Infos
		expected expected
	}{
		{
			desc:  "all routers, but no config",
			path:  "/http/routers",
			infos: &runtime.Infos{},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "1",
				jsonFile:   "testdata/routers-empty.json",
			},
		},
		{
			desc: "all routers",
			path: "/http/routers",
			infos: &runtime.Infos{
				Routers: map[string]*runtime.RouterInfo{
					"test@myprovider": {
						Router: &dynamic.Router{
							Router: &traefikdynamic.Router{
								EntryPoints: []string{"web"},
								Service:     "foo-service@myprovider",
								Rule:        "Host(`foo.bar.other`)",
								Middlewares: []string{"addPrefixTest", "auth"},
							},
						},
					},
					"bar@myprovider": {
						Router: &dynamic.Router{
							Router: &traefikdynamic.Router{
								EntryPoints: []string{"web"},
								Service:     "foo-service@myprovider",
								Rule:        "Host(`foo.bar`)",
								Middlewares: []string{"auth", "addPrefixTest@anotherprovider"},
							},
						},
					},
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "1",
				jsonFile:   "testdata/routers.json",
			},
		},
		{
			desc: "all routers, pagination, 1 res per page, want page 2",
			path: "/http/routers?page=2&per_page=1",
			infos: &runtime.Infos{
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
					},
					"baz@myprovider": {
						Router: &dynamic.Router{
							Router: &traefikdynamic.Router{
								EntryPoints: []string{"web"},
								Service:     "foo-service@myprovider",
								Rule:        "Host(`toto.bar`)",
							},
						},
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
					},
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "3",
				jsonFile:   "testdata/routers-page2.json",
			},
		},
		{
			desc: "all routers, pagination, 19 results overall, 7 res per page, want page 3",
			path: "/http/routers?page=3&per_page=7",
			infos: &runtime.Infos{
				Routers: generateHTTPRouters(19),
			},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "1",
				jsonFile:   "testdata/routers-many-lastpage.json",
			},
		},
		{
			desc: "all routers, pagination, 5 results overall, 10 res per page, want page 2",
			path: "/http/routers?page=2&per_page=10",
			infos: &runtime.Infos{
				Routers: generateHTTPRouters(5),
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			desc: "all routers, pagination, 10 results overall, 10 res per page, want page 2",
			path: "/http/routers?page=2&per_page=10",
			infos: &runtime.Infos{
				Routers: generateHTTPRouters(10),
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			desc: "routers filtered by status",
			path: "/http/routers?status=enabled",
			infos: &runtime.Infos{
				Routers: map[string]*runtime.RouterInfo{
					"test@myprovider": {
						Router: &dynamic.Router{
							Router: &traefikdynamic.Router{
								EntryPoints: []string{"web"},
								Service:     "foo-service@myprovider",
								Rule:        "Host(`foo.bar.other`)",
								Middlewares: []string{"addPrefixTest", "auth"},
							},
						},
						Status: runtime.StatusEnabled,
					},
					"bar@myprovider": {
						Router: &dynamic.Router{
							Router: &traefikdynamic.Router{
								EntryPoints: []string{"web"},
								Service:     "foo-service@myprovider",
								Rule:        "Host(`foo.bar`)",
								Middlewares: []string{"auth", "addPrefixTest@anotherprovider"},
							},
						},
						Status: runtime.StatusDisabled,
					},
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "1",
				jsonFile:   "testdata/routers-filtered-status.json",
			},
		},
		{
			desc: "routers filtered by search",
			path: "/http/routers?search=fii",
			infos: &runtime.Infos{
				Routers: map[string]*runtime.RouterInfo{
					"test@myprovider": {
						Router: &dynamic.Router{
							Router: &traefikdynamic.Router{
								EntryPoints: []string{"web"},
								Service:     "fii-service@myprovider",
								Rule:        "Host(`fii.bar.other`)",
								Middlewares: []string{"addPrefixTest", "auth"},
							},
						},
						Status: runtime.StatusEnabled,
					},
					"bar@myprovider": {
						Router: &dynamic.Router{
							Router: &traefikdynamic.Router{
								EntryPoints: []string{"web"},
								Service:     "foo-service@myprovider",
								Rule:        "Host(`foo.bar`)",
								Middlewares: []string{"auth", "addPrefixTest@anotherprovider"},
							},
						},
						Status: runtime.StatusDisabled,
					},
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "1",
				jsonFile:   "testdata/routers-filtered-search.json",
			},
		},
		{
			desc: "one router by id",
			path: "/http/routers/bar@myprovider",
			infos: &runtime.Infos{
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
						Status: "enabled",
					},
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				jsonFile:   "testdata/router-bar.json",
			},
		},
		{
			desc: "one router by id, that does not exist",
			path: "/http/routers/foo@myprovider",
			infos: &runtime.Infos{
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
					},
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
		{
			desc:  "one router by id, but no config",
			path:  "/http/routers/foo@myprovider",
			infos: &runtime.Infos{},
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
		{
			desc:  "all services, but no config",
			path:  "/http/services",
			infos: &runtime.Infos{},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "1",
				jsonFile:   "testdata/services-empty.json",
			},
		},
		{
			desc: "all services",
			path: "/http/services",
			infos: &runtime.Infos{
				Services: map[string]*runtime.ServiceInfo{
					"bar@myprovider": func() *runtime.ServiceInfo {
						si := &runtime.ServiceInfo{
							Service: &dynamic.Service{
								ReverseProxy: &dynamic.ReverseProxy{
									PassHostHeader: utils.Bool(true),
									LoadBalancer: &dynamic.LoadBalancer{
										Servers: []*dynamic.Server{
											{
												URL: "http://127.0.0.1",
											},
										},
									},
								},
							},
							UsedBy: []string{"foo@myprovider", "test@myprovider"},
						}
						si.UpdateServerStatus("http://127.0.0.1", "UP")
						return si
					}(),
					"baz@myprovider": func() *runtime.ServiceInfo {
						si := &runtime.ServiceInfo{
							Service: &dynamic.Service{
								ReverseProxy: &dynamic.ReverseProxy{
									PassHostHeader: utils.Bool(true),
									LoadBalancer: &dynamic.LoadBalancer{
										Servers: []*dynamic.Server{
											{
												URL: "http://127.0.0.2",
											},
										},
										Sticky: &dynamic.Sticky{
											Cookie: &dynamic.Cookie{
												Name:     "chocolat",
												Secure:   true,
												HTTPOnly: true,
											},
										},
									},
								},
							},
							UsedBy: []string{"foo@myprovider"},
						}
						si.UpdateServerStatus("http://127.0.0.2", "UP")
						return si
					}(),
					"mock@myprovider": {
						Service: &dynamic.Service{
							Mock: &dynamic.Mock{},
						},
						Status: runtime.StatusEnabled,
						UsedBy: []string{"foo@myprovider"},
					},
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "1",
				jsonFile:   "testdata/services.json",
			},
		},
		{
			desc: "all services, 1 res per page, want page 2",
			path: "/http/services?page=2&per_page=1",
			infos: &runtime.Infos{
				Services: map[string]*runtime.ServiceInfo{
					"bar@myprovider": func() *runtime.ServiceInfo {
						si := &runtime.ServiceInfo{
							Service: &dynamic.Service{
								ReverseProxy: &dynamic.ReverseProxy{
									PassHostHeader: utils.Bool(true),
									LoadBalancer: &dynamic.LoadBalancer{
										Servers: []*dynamic.Server{
											{
												URL: "http://127.0.0.1",
											},
										},
									},
								},
							},
							UsedBy: []string{"foo@myprovider", "test@myprovider"},
						}
						si.UpdateServerStatus("http://127.0.0.1", "UP")
						return si
					}(),
					"baz@myprovider": func() *runtime.ServiceInfo {
						si := &runtime.ServiceInfo{
							Service: &dynamic.Service{
								ReverseProxy: &dynamic.ReverseProxy{
									PassHostHeader: utils.Bool(true),
									LoadBalancer: &dynamic.LoadBalancer{
										Servers: []*dynamic.Server{
											{
												URL: "http://127.0.0.2",
											},
										},
									},
								},
							},
							UsedBy: []string{"foo@myprovider"},
						}
						si.UpdateServerStatus("http://127.0.0.2", "UP")
						return si
					}(),
					"test@myprovider": func() *runtime.ServiceInfo {
						si := &runtime.ServiceInfo{
							Service: &dynamic.Service{
								ReverseProxy: &dynamic.ReverseProxy{
									PassHostHeader: utils.Bool(true),
									LoadBalancer: &dynamic.LoadBalancer{
										Servers: []*dynamic.Server{
											{
												URL: "http://127.0.0.3",
											},
										},
									},
								},
							},
							UsedBy: []string{"foo@myprovider", "test@myprovider"},
						}
						si.UpdateServerStatus("http://127.0.0.4", "UP")
						return si
					}(),
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "3",
				jsonFile:   "testdata/services-page2.json",
			},
		},
		{
			desc: "services filtered by status",
			path: "/http/services?status=enabled",
			infos: &runtime.Infos{
				Services: map[string]*runtime.ServiceInfo{
					"bar@myprovider": func() *runtime.ServiceInfo {
						si := &runtime.ServiceInfo{
							Service: &dynamic.Service{
								ReverseProxy: &dynamic.ReverseProxy{
									PassHostHeader: utils.Bool(true),
									LoadBalancer: &dynamic.LoadBalancer{
										Servers: []*dynamic.Server{
											{
												URL: "http://127.0.0.1",
											},
										},
									},
								},
							},
							UsedBy: []string{"foo@myprovider", "test@myprovider"},
							Status: runtime.StatusEnabled,
						}
						si.UpdateServerStatus("http://127.0.0.1", "UP")
						return si
					}(),
					"baz@myprovider": func() *runtime.ServiceInfo {
						si := &runtime.ServiceInfo{
							Service: &dynamic.Service{
								ReverseProxy: &dynamic.ReverseProxy{
									PassHostHeader: utils.Bool(true),
									LoadBalancer: &dynamic.LoadBalancer{
										Servers: []*dynamic.Server{
											{
												URL: "http://127.0.0.2",
											},
										},
									},
								},
							},
							UsedBy: []string{"foo@myprovider"},
							Status: runtime.StatusDisabled,
						}
						si.UpdateServerStatus("http://127.0.0.2", "UP")
						return si
					}(),
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "1",
				jsonFile:   "testdata/services-filtered-status.json",
			},
		},
		{
			desc: "services filtered by search",
			path: "/http/services?search=baz",
			infos: &runtime.Infos{
				Services: map[string]*runtime.ServiceInfo{
					"bar@myprovider": func() *runtime.ServiceInfo {
						si := &runtime.ServiceInfo{
							Service: &dynamic.Service{
								ReverseProxy: &dynamic.ReverseProxy{
									PassHostHeader: utils.Bool(true),
									LoadBalancer: &dynamic.LoadBalancer{
										Servers: []*dynamic.Server{
											{
												URL: "http://127.0.0.1",
											},
										},
									},
								},
							},
							UsedBy: []string{"foo@myprovider", "test@myprovider"},
							Status: runtime.StatusEnabled,
						}
						si.UpdateServerStatus("http://127.0.0.1", "UP")
						return si
					}(),
					"baz@myprovider": func() *runtime.ServiceInfo {
						si := &runtime.ServiceInfo{
							Service: &dynamic.Service{
								ReverseProxy: &dynamic.ReverseProxy{
									PassHostHeader: utils.Bool(true),
									LoadBalancer: &dynamic.LoadBalancer{
										Servers: []*dynamic.Server{
											{
												URL: "http://127.0.0.2",
											},
										},
									},
								},
							},
							UsedBy: []string{"foo@myprovider"},
							Status: runtime.StatusDisabled,
						}
						si.UpdateServerStatus("http://127.0.0.2", "UP")
						return si
					}(),
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "1",
				jsonFile:   "testdata/services-filtered-search.json",
			},
		},
		{
			desc: "one service by id",
			path: "/http/services/bar@myprovider",
			infos: &runtime.Infos{
				Services: map[string]*runtime.ServiceInfo{
					"bar@myprovider": func() *runtime.ServiceInfo {
						si := &runtime.ServiceInfo{
							Service: &dynamic.Service{
								ReverseProxy: &dynamic.ReverseProxy{
									PassHostHeader: utils.Bool(true),
									LoadBalancer: &dynamic.LoadBalancer{
										Servers: []*dynamic.Server{
											{
												URL: "http://127.0.0.1",
											},
										},
									},
								},
							},
							UsedBy: []string{"foo@myprovider", "test@myprovider"},
						}
						si.UpdateServerStatus("http://127.0.0.1", "UP")
						return si
					}(),
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				jsonFile:   "testdata/service-bar.json",
			},
		},
		{
			desc: "one service by id, that does not exist",
			path: "/http/services/nono@myprovider",
			infos: &runtime.Infos{
				Services: map[string]*runtime.ServiceInfo{
					"bar@myprovider": func() *runtime.ServiceInfo {
						si := &runtime.ServiceInfo{
							Service: &dynamic.Service{
								ReverseProxy: &dynamic.ReverseProxy{
									PassHostHeader: utils.Bool(true),
									LoadBalancer: &dynamic.LoadBalancer{
										Servers: []*dynamic.Server{
											{
												URL: "http://127.0.0.1",
											},
										},
									},
								},
							},
							UsedBy: []string{"foo@myprovider", "test@myprovider"},
						}
						si.UpdateServerStatus("http://127.0.0.1", "UP")
						return si
					}(),
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
		{
			desc:  "one service by id, but no config",
			path:  "/http/services/foo@myprovider",
			infos: &runtime.Infos{},
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
		{
			desc:  "all middlewares, but no config",
			path:  "/http/middlewares",
			infos: &runtime.Infos{},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "1",
				jsonFile:   "testdata/middlewares-empty.json",
			},
		},
		{
			desc: "all middlewares",
			path: "/http/middlewares",
			infos: &runtime.Infos{
				Middlewares: map[string]*runtime.MiddlewareInfo{
					"auth@myprovider": {
						Middleware: &dynamic.Middleware{
							Auth: &dynamic.Auth{},
						},
						UsedBy: []string{"bar@myprovider", "test@myprovider"},
					},
					"addPrefixTest@myprovider": {
						Middleware: &dynamic.Middleware{
							Middleware: &traefikdynamic.Middleware{
								AddPrefix: &traefikdynamic.AddPrefix{
									Prefix: "/titi",
								},
							},
						},
						UsedBy: []string{"test@myprovider"},
					},
					"addPrefixTest@anotherprovider": {
						Middleware: &dynamic.Middleware{
							Middleware: &traefikdynamic.Middleware{
								AddPrefix: &traefikdynamic.AddPrefix{
									Prefix: "/toto",
								},
							},
						},
						UsedBy: []string{"bar@myprovider"},
					},
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "1",
				jsonFile:   "testdata/middlewares.json",
			},
		},
		{
			desc: "middlewares filtered by status",
			path: "/http/middlewares?status=enabled",
			infos: &runtime.Infos{
				Middlewares: map[string]*runtime.MiddlewareInfo{
					"auth@myprovider": {
						Middleware: &dynamic.Middleware{
							Auth: &dynamic.Auth{},
						},
						UsedBy: []string{"bar@myprovider", "test@myprovider"},
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
						UsedBy: []string{"test@myprovider"},
						Status: runtime.StatusDisabled,
					},
					"addPrefixTest@anotherprovider": {
						Middleware: &dynamic.Middleware{
							Middleware: &traefikdynamic.Middleware{
								AddPrefix: &traefikdynamic.AddPrefix{
									Prefix: "/toto",
								},
							},
						},
						UsedBy: []string{"bar@myprovider"},
						Status: runtime.StatusEnabled,
					},
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "1",
				jsonFile:   "testdata/middlewares-filtered-status.json",
			},
		},
		{
			desc: "middlewares filtered by search",
			path: "/http/middlewares?search=addprefixtest",
			infos: &runtime.Infos{
				Middlewares: map[string]*runtime.MiddlewareInfo{
					"auth@myprovider": {
						Middleware: &dynamic.Middleware{
							Auth: &dynamic.Auth{},
						},
						UsedBy: []string{"bar@myprovider", "test@myprovider"},
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
						UsedBy: []string{"test@myprovider"},
						Status: runtime.StatusDisabled,
					},
					"addPrefixTest@anotherprovider": {
						Middleware: &dynamic.Middleware{
							Middleware: &traefikdynamic.Middleware{
								AddPrefix: &traefikdynamic.AddPrefix{
									Prefix: "/toto",
								},
							},
						},
						UsedBy: []string{"bar@myprovider"},
						Status: runtime.StatusEnabled,
					},
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "1",
				jsonFile:   "testdata/middlewares-filtered-search.json",
			},
		},
		{
			desc: "all middlewares, 1 res per page, want page 2",
			path: "/http/middlewares?page=2&per_page=1",
			infos: &runtime.Infos{
				Middlewares: map[string]*runtime.MiddlewareInfo{
					"auth@myprovider": {
						Middleware: &dynamic.Middleware{
							Auth: &dynamic.Auth{},
						},
						UsedBy: []string{"bar@myprovider", "test@myprovider"},
					},
					"addPrefixTest@myprovider": {
						Middleware: &dynamic.Middleware{
							Middleware: &traefikdynamic.Middleware{
								AddPrefix: &traefikdynamic.AddPrefix{
									Prefix: "/titi",
								},
							},
						},
						UsedBy: []string{"test@myprovider"},
					},
					"addPrefixTest@anotherprovider": {
						Middleware: &dynamic.Middleware{
							Middleware: &traefikdynamic.Middleware{
								AddPrefix: &traefikdynamic.AddPrefix{
									Prefix: "/toto",
								},
							},
						},
						UsedBy: []string{"bar@myprovider"},
					},
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				nextPage:   "3",
				jsonFile:   "testdata/middlewares-page2.json",
			},
		},
		{
			desc: "one middleware by id",
			path: "/http/middlewares/auth@myprovider",
			infos: &runtime.Infos{
				Middlewares: map[string]*runtime.MiddlewareInfo{
					"auth@myprovider": {
						Middleware: &dynamic.Middleware{
							Auth: &dynamic.Auth{},
						},
						UsedBy: []string{"bar@myprovider", "test@myprovider"},
					},
					"addPrefixTest@myprovider": {
						Middleware: &dynamic.Middleware{
							Middleware: &traefikdynamic.Middleware{
								AddPrefix: &traefikdynamic.AddPrefix{
									Prefix: "/titi",
								},
							},
						},
						UsedBy: []string{"test@myprovider"},
					},
					"addPrefixTest@anotherprovider": {
						Middleware: &dynamic.Middleware{
							Middleware: &traefikdynamic.Middleware{
								AddPrefix: &traefikdynamic.AddPrefix{
									Prefix: "/toto",
								},
							},
						},
						UsedBy: []string{"bar@myprovider"},
					},
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				jsonFile:   "testdata/middleware-auth.json",
			},
		},
		{
			desc: "one middleware by id, that does not exist",
			path: "/http/middlewares/foo@myprovider",
			infos: &runtime.Infos{
				Middlewares: map[string]*runtime.MiddlewareInfo{
					"auth@myprovider": {
						Middleware: &dynamic.Middleware{
							Auth: &dynamic.Auth{},
						},
						UsedBy: []string{"bar@myprovider", "test@myprovider"},
					},
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
		{
			desc:  "one middleware by id, but no config",
			path:  "/http/middlewares/foo@myprovider",
			infos: &runtime.Infos{},
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
		{
			desc: "Get rawdata",
			path: "/rawdata",
			infos: &runtime.Infos{
				Services: map[string]*runtime.ServiceInfo{
					"foo-service@myprovider": {
						Service: &dynamic.Service{
							ReverseProxy: &dynamic.ReverseProxy{
								PassHostHeader: utils.Bool(true),
								LoadBalancer: &dynamic.LoadBalancer{
									Servers: []*dynamic.Server{
										{
											URL: "http://127.0.0.1",
										},
									},
								},
							},
						},
					},
				},
				Middlewares: map[string]*runtime.MiddlewareInfo{
					"auth@myprovider": {
						Middleware: &dynamic.Middleware{
							Auth: &dynamic.Auth{},
						},
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
					},
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				jsonFile:   "testdata/getrawdata.json",
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			// To lazily initialize the Statuses.
			test.infos.PopulateUsedBy()
			test.infos.RouterInfosByEntryPoint(context.Background(), []string{"web"}, false)

			router := mux.NewRouter()
			handler := NewHTTP(test.infos)
			handler.Append(router)
			server := httptest.NewServer(router)

			resp, err := http.DefaultClient.Get(server.URL + test.path)
			require.NoError(t, err)

			require.Equal(t, test.expected.statusCode, resp.StatusCode)

			assert.Equal(t, test.expected.nextPage, resp.Header.Get(nextPageHeader))

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

func generateHTTPRouters(nbRouters int) map[string]*runtime.RouterInfo {
	routers := make(map[string]*runtime.RouterInfo, nbRouters)
	for i := 0; i < nbRouters; i++ {
		routers[fmt.Sprintf("bar%2d@myprovider", i)] = &runtime.RouterInfo{
			Router: &dynamic.Router{
				Router: &traefikdynamic.Router{
					EntryPoints: []string{"web"},
					Service:     "foo-service@myprovider",
					Rule:        "Host(`foo.bar" + strconv.Itoa(i) + "`)",
				},
			},
		}
	}
	return routers
}
