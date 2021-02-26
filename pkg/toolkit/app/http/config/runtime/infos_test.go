// +build unit

package runtime

import (
	"context"
	"testing"

	testutilsdynamic "github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/config/dynamic/testutils"
	"github.com/stretchr/testify/assert"
)

func TestInfos(t *testing.T) {
	infos := NewInfos(testutilsdynamic.Config)
	infos.PopulateUsedBy()

	assert.Len(t, infos.Routers, 2, "Count of routers")
	assert.Len(t, infos.Middlewares, 3, "Count of middlewares")
	assert.Len(t, infos.Services, 2, "Count of services")

	assert.Equal(t, infos.Routers["router-proxy@provider1"].Router, testutilsdynamic.Config.HTTP.Routers["router-proxy@provider1"], "Router #1")
	assert.Equal(t, infos.Routers["router-dashboard@provider2"].Router, testutilsdynamic.Config.HTTP.Routers["router-dashboard@provider2"], "Router #2")
	assert.Equal(t, infos.Middlewares["middleware-foo@provider1"].Middleware, testutilsdynamic.Config.HTTP.Middlewares["middleware-foo@provider1"], "Middleware #1")
	assert.Equal(t, infos.Middlewares["middleware-bar@provider1"].Middleware, testutilsdynamic.Config.HTTP.Middlewares["middleware-bar@provider1"], "Middleware #2")
	assert.Equal(t, infos.Middlewares["accesslog@provider2"].Middleware, testutilsdynamic.Config.HTTP.Middlewares["accesslog@provider2"], "Middleware #3")
	assert.Equal(t, infos.Services["proxy@provider1"].Service, testutilsdynamic.Config.HTTP.Services["proxy@provider1"], "Service #1")
	assert.Equal(t, infos.Services["dashboard@provider1"].Service, testutilsdynamic.Config.HTTP.Services["dashboard@provider1"], "Service #2")

	assert.Equal(t, []string{"router-proxy@provider1"}, infos.Middlewares["middleware-foo@provider1"].UsedBy, "#1 UsedBy")
	assert.Equal(t, []string{"router-dashboard@provider2"}, infos.Services["dashboard@provider1"].UsedBy, "#2 UsedBy")

	rtInfos := infos.RouterInfosByEntryPoint(context.Background(), []string{"ep-foo", "ep-bar"}, false)
	assert.Len(t, rtInfos, 2)
	assert.Len(t, rtInfos["ep-foo"], 2)
	assert.Contains(t, rtInfos["ep-foo"], "router-proxy@provider1")
	assert.Contains(t, rtInfos["ep-foo"], "router-dashboard@provider2")
	assert.Len(t, rtInfos["ep-bar"], 1)
	assert.Contains(t, rtInfos["ep-bar"], "router-proxy@provider1")
}
