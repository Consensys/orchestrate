// +build unit

package dashboard
// 
// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"strconv"
// 	"testing"
// 
// 	traefikstatic "github.com/traefik/traefik/v2/pkg/config/static"
// 	traefiktypes "github.com/traefik/paerser/types"
// 	"github.com/gorilla/mux"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )
// 
// func TestEntryPoints(t *testing.T) {
// 	type expected struct {
// 		statusCode int
// 		nextPage   string
// 		jsonFile   string
// 	}
// 
// 	testCases := []struct {
// 		desc     string
// 		path     string
// 		conf     *traefikstatic.Configuration
// 		expected expected
// 	}{
// 		{
// 			desc: "all entry points, but no config",
// 			path: "/entrypoints",
// 			conf: &traefikstatic.Configuration{API: &traefikstatic.API{}, Global: &traefikstatic.Global{}},
// 			expected: expected{
// 				statusCode: http.StatusOK,
// 				nextPage:   "1",
// 				jsonFile:   "testdata/entrypoints-empty.json",
// 			},
// 		},
// 		{
// 			desc: "all entry points",
// 			path: "/entrypoints",
// 			conf: &traefikstatic.Configuration{
// 				Global: &traefikstatic.Global{},
// 				API:    &traefikstatic.API{},
// 				EntryPoints: map[string]*traefikstatic.EntryPoint{
// 					"web": {
// 						Address: ":80",
// 						Transport: &traefikstatic.EntryPointsTransport{
// 							LifeCycle: &traefikstatic.LifeCycle{
// 								RequestAcceptGraceTimeout: traefiktypes.Duration(1),
// 								GraceTimeOut:              traefiktypes.Duration(2),
// 							},
// 							RespondingTimeouts: &traefikstatic.RespondingTimeouts{
// 								ReadTimeout:  3,
// 								WriteTimeout: 4,
// 								IdleTimeout:  5,
// 							},
// 						},
// 						ProxyProtocol: &traefikstatic.ProxyProtocol{
// 							Insecure:   true,
// 							TrustedIPs: []string{"192.168.1.1", "192.168.1.2"},
// 						},
// 						ForwardedHeaders: &traefikstatic.ForwardedHeaders{
// 							Insecure:   true,
// 							TrustedIPs: []string{"192.168.1.3", "192.168.1.4"},
// 						},
// 					},
// 					"websecure": {
// 						Address: ":443",
// 						Transport: &traefikstatic.EntryPointsTransport{
// 							LifeCycle: &traefikstatic.LifeCycle{
// 								RequestAcceptGraceTimeout: 10,
// 								GraceTimeOut:              20,
// 							},
// 							RespondingTimeouts: &traefikstatic.RespondingTimeouts{
// 								ReadTimeout:  30,
// 								WriteTimeout: 40,
// 								IdleTimeout:  50,
// 							},
// 						},
// 						ProxyProtocol: &traefikstatic.ProxyProtocol{
// 							Insecure:   true,
// 							TrustedIPs: []string{"192.168.1.10", "192.168.1.20"},
// 						},
// 						ForwardedHeaders: &traefikstatic.ForwardedHeaders{
// 							Insecure:   true,
// 							TrustedIPs: []string{"192.168.1.30", "192.168.1.40"},
// 						},
// 					},
// 				},
// 			},
// 			expected: expected{
// 				statusCode: http.StatusOK,
// 				nextPage:   "1",
// 				jsonFile:   "testdata/entrypoints.json",
// 			},
// 		},
// 		{
// 			desc: "all entry points, pagination, 1 res per page, want page 2",
// 			path: "/entrypoints?page=2&per_page=1",
// 			conf: &traefikstatic.Configuration{
// 				Global: &traefikstatic.Global{},
// 				API:    &traefikstatic.API{},
// 				EntryPoints: map[string]*traefikstatic.EntryPoint{
// 					"web1": {Address: ":81"},
// 					"web2": {Address: ":82"},
// 					"web3": {Address: ":83"},
// 				},
// 			},
// 			expected: expected{
// 				statusCode: http.StatusOK,
// 				nextPage:   "3",
// 				jsonFile:   "testdata/entrypoints-page2.json",
// 			},
// 		},
// 		{
// 			desc: "all entry points, pagination, 19 results overall, 7 res per page, want page 3",
// 			path: "/entrypoints?page=3&per_page=7",
// 			conf: &traefikstatic.Configuration{
// 				Global:      &traefikstatic.Global{},
// 				API:         &traefikstatic.API{},
// 				EntryPoints: generateEntryPoints(19),
// 			},
// 			expected: expected{
// 				statusCode: http.StatusOK,
// 				nextPage:   "1",
// 				jsonFile:   "testdata/entrypoints-many-lastpage.json",
// 			},
// 		},
// 		{
// 			desc: "all entry points, pagination, 5 results overall, 10 res per page, want page 2",
// 			path: "/entrypoints?page=2&per_page=10",
// 			conf: &traefikstatic.Configuration{
// 				Global:      &traefikstatic.Global{},
// 				API:         &traefikstatic.API{},
// 				EntryPoints: generateEntryPoints(5),
// 			},
// 			expected: expected{
// 				statusCode: http.StatusBadRequest,
// 			},
// 		},
// 		{
// 			desc: "all entry points, pagination, 10 results overall, 10 res per page, want page 2",
// 			path: "/entrypoints?page=2&per_page=10",
// 			conf: &traefikstatic.Configuration{
// 				Global:      &traefikstatic.Global{},
// 				API:         &traefikstatic.API{},
// 				EntryPoints: generateEntryPoints(10),
// 			},
// 			expected: expected{
// 				statusCode: http.StatusBadRequest,
// 			},
// 		},
// 		{
// 			desc: "one entry point by id",
// 			path: "/entrypoints/bar",
// 			conf: &traefikstatic.Configuration{
// 				Global: &traefikstatic.Global{},
// 				API:    &traefikstatic.API{},
// 				EntryPoints: map[string]*traefikstatic.EntryPoint{
// 					"bar": {Address: ":81"},
// 				},
// 			},
// 			expected: expected{
// 				statusCode: http.StatusOK,
// 				jsonFile:   "testdata/entrypoint-bar.json",
// 			},
// 		},
// 		{
// 			desc: "one entry point by id, that does not exist",
// 			path: "/entrypoints/foo",
// 			conf: &traefikstatic.Configuration{
// 				Global: &traefikstatic.Global{},
// 				API:    &traefikstatic.API{},
// 				EntryPoints: map[string]*traefikstatic.EntryPoint{
// 					"bar": {Address: ":81"},
// 				},
// 			},
// 			expected: expected{
// 				statusCode: http.StatusNotFound,
// 			},
// 		},
// 		{
// 			desc: "one entry point by id, but no config",
// 			path: "/entrypoints/foo",
// 			conf: &traefikstatic.Configuration{API: &traefikstatic.API{}, Global: &traefikstatic.Global{}},
// 			expected: expected{
// 				statusCode: http.StatusNotFound,
// 			},
// 		},
// 	}
// 
// 	for _, test := range testCases {
// 		test := test
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Parallel()
// 
// 			router := mux.NewRouter()
// 			handler := NewEntryPoint(test.conf)
// 			handler.Append(router)
// 			server := httptest.NewServer(router)
// 
// 			resp, err := http.DefaultClient.Get(server.URL + test.path)
// 			require.NoError(t, err)
// 
// 			require.Equal(t, test.expected.statusCode, resp.StatusCode)
// 
// 			assert.Equal(t, test.expected.nextPage, resp.Header.Get(nextPageHeader))
// 
// 			if test.expected.jsonFile == "" {
// 				return
// 			}
// 
// 			assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
// 			contents, err := ioutil.ReadAll(resp.Body)
// 			require.NoError(t, err)
// 
// 			err = resp.Body.Close()
// 			require.NoError(t, err)
// 
// 			if *updateExpected {
// 				var results interface{}
// 				err = json.Unmarshal(contents, &results)
// 				require.NoError(t, err)
// 
// 				newJSON, e := json.MarshalIndent(results, "", "\t")
// 				require.NoError(t, e)
// 
// 				err = ioutil.WriteFile(test.expected.jsonFile, newJSON, 0644)
// 				require.NoError(t, err)
// 			}
// 
// 			data, err := ioutil.ReadFile(test.expected.jsonFile)
// 			require.NoError(t, err)
// 			assert.JSONEq(t, string(data), string(contents))
// 		})
// 	}
// }
// 
// func generateEntryPoints(nb int) map[string]*traefikstatic.EntryPoint {
// 	eps := make(map[string]*traefikstatic.EntryPoint, nb)
// 	for i := 0; i < nb; i++ {
// 		eps[fmt.Sprintf("ep%2d", i)] = &traefikstatic.EntryPoint{
// 			Address: ":" + strconv.Itoa(i),
// 		}
// 	}
// 
// 	return eps
// }
