package runtime

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

// Status of the router/middleware/service
const (
	StatusEnabled  = "enabled"
	StatusDisabled = "disabled"
	StatusWarning  = "warning"
)

// Infos holds the information about the currently running traefik instance.
type Infos struct {
	Routers     map[string]*RouterInfo     `json:"routers,omitempty"`
	Middlewares map[string]*MiddlewareInfo `json:"middlewares,omitempty"`
	Services    map[string]*ServiceInfo    `json:"services,omitempty"`
}

// NewInfos returns Infos initialized with the given conf. It never returns nil.
func NewInfos(conf *dynamic.Configuration) *Infos {
	if conf.HTTP == nil {
		return &Infos{}
	}

	return &Infos{
		Routers:     NewRouterInfos(conf.HTTP.Routers),
		Middlewares: NewMiddlewareInfos(conf.HTTP.Middlewares),
		Services:    NewServiceInfos(conf.HTTP.Services),
	}
}

func NewRouterInfos(routers map[string]*dynamic.Router) map[string]*RouterInfo {
	if len(routers) > 0 {
		infos := make(map[string]*RouterInfo, len(routers))
		for k, v := range routers {
			infos[k] = &RouterInfo{Router: v, Status: StatusEnabled}
		}
		return infos
	}
	return nil
}

func NewMiddlewareInfos(middlewares map[string]*dynamic.Middleware) map[string]*MiddlewareInfo {
	if len(middlewares) > 0 {
		infos := make(map[string]*MiddlewareInfo, len(middlewares))
		for k, v := range middlewares {
			infos[k] = &MiddlewareInfo{Middleware: v, Status: StatusEnabled}
		}
		return infos

	}
	return nil
}

func NewServiceInfos(services map[string]*dynamic.Service) map[string]*ServiceInfo {
	if len(services) > 0 {
		infos := make(map[string]*ServiceInfo, len(services))
		for k, v := range services {
			infos[k] = &ServiceInfo{Service: v, Status: StatusEnabled}
		}
		return infos

	}
	return nil
}

// PopulateUsedBy populates all the UsedBy lists of the underlying fields of r,
// based on the relations between the included services, routers, and middlewares.
func (infos *Infos) PopulateUsedBy() {
	if infos == nil {
		return
	}

	for routerName, routerInfo := range infos.Routers {
		// lazily initialize Status in case caller forgot to do it
		if routerInfo.Status == "" {
			routerInfo.Status = StatusEnabled
		}

		providerName := provider.GetName(routerName)
		if providerName == "" {
			continue
		}

		for _, midName := range routerInfo.Router.Middlewares {
			qualifiedMidName := provider.QualifyName(providerName, midName)
			if _, ok := infos.Middlewares[qualifiedMidName]; !ok {
				continue
			}
			infos.Middlewares[qualifiedMidName].UsedBy = append(infos.Middlewares[qualifiedMidName].UsedBy, routerName)
		}

		qualifiedMidNameSrvName := provider.QualifyName(providerName, routerInfo.Router.Service)
		if _, ok := infos.Services[qualifiedMidNameSrvName]; !ok {
			continue
		}

		infos.Services[qualifiedMidNameSrvName].UsedBy = append(infos.Services[qualifiedMidNameSrvName].UsedBy, routerName)
	}

	for k, serviceInfo := range infos.Services {
		// lazily initialize Status in case caller forgot to do it
		if serviceInfo.Status == "" {
			serviceInfo.Status = StatusEnabled
		}

		sort.Strings(infos.Services[k].UsedBy)
	}

	for midName, mid := range infos.Middlewares {
		// lazily initialize Status in case caller forgot to do it
		if mid.Status == "" {
			mid.Status = StatusEnabled
		}

		sort.Strings(infos.Middlewares[midName].UsedBy)
	}
}

// GetRouters returns all routers infos by entry points name and routers name.
func (infos *Infos) RouterInfosByEntryPoint(ctx context.Context, entryPointNames []string, tls bool) map[string]map[string]*RouterInfo {
	rtInfos := make(map[string]map[string]*RouterInfo)

	for rtName, rt := range infos.Routers {
		if (tls && rt.TLS == nil) || (!tls && rt.TLS != nil) {
			continue
		}

		logger := log.FromContext(log.With(ctx, log.Str(log.RouterName, rtName)))

		eps := rt.EntryPoints
		if len(eps) == 0 {
			eps = entryPointNames
		}

		entryPointsCount := 0
		for _, entryPointName := range eps {
			if !utils.ContainsString(entryPointNames, entryPointName) {
				rt.AddError(fmt.Errorf("entryPoint %q doesn't exist", entryPointName), false)
				logger.WithField(log.EntryPointName, entryPointName).
					Errorf("entryPoint %q doesn't exist", entryPointName)
				continue
			}

			if _, ok := rtInfos[entryPointName]; !ok {
				rtInfos[entryPointName] = make(map[string]*RouterInfo)
			}

			entryPointsCount++
			rt.Using = append(rt.Using, entryPointName)

			rtInfos[entryPointName][rtName] = rt
		}

		if entryPointsCount == 0 {
			rt.AddError(fmt.Errorf("no valid entryPoint for this router"), true)
			logger.Error("no valid entryPoint for this router")
		}

		rt.Using = utils.UniqueString(rt.Using)
	}

	return rtInfos
}

type MiddlewareInfo struct {
	*dynamic.Middleware // dynamic configuration
	// Err contains all the errors that occurred during service creation.
	Err    []string `json:"error,omitempty"`
	Status string   `json:"status,omitempty"`
	UsedBy []string `json:"usedBy,omitempty"` // list of routers and services using that middleware.
}

// AddError adds err to s.Err, if it does not already exist.
// If critical is set, m is marked as disabled.
func (m *MiddlewareInfo) AddError(err error, critical bool) {
	for _, value := range m.Err {
		if value == err.Error() {
			return
		}
	}

	m.Err = append(m.Err, err.Error())
	if critical {
		m.Status = StatusDisabled
		return
	}

	// only set it to "warning" if not already in a worse state
	if m.Status != StatusDisabled {
		m.Status = StatusWarning
	}
}

// ServiceInfo holds information about a currently running service.
type ServiceInfo struct {
	*dynamic.Service // dynamic configuration
	// Err contains all the errors that occurred during service creation.
	Err []string `json:"error,omitempty"`
	// Status reports whether the service is disabled, in a warning state, or all good (enabled).
	// If not in "enabled" state, the reason for it should be in the list of Err.
	// It is the caller's responsibility to set the initial status.
	Status string   `json:"status,omitempty"`
	UsedBy []string `json:"usedBy,omitempty"` // list of routers using that service

	serverStatusMu sync.RWMutex
	serverStatus   map[string]string // keyed by server URL
}

// AddError adds err to s.Err, if it does not already exist.
// If critical is set, s is marked as disabled.
func (s *ServiceInfo) AddError(err error, critical bool) {
	for _, value := range s.Err {
		if value == err.Error() {
			return
		}
	}

	s.Err = append(s.Err, err.Error())
	if critical {
		s.Status = StatusDisabled
		return
	}

	// only set it to "warning" if not already in a worse state
	if s.Status != StatusDisabled {
		s.Status = StatusWarning
	}
}

// UpdateServerStatus sets the status of the server in the ServiceInfo.
// It is the responsibility of the caller to check that s is not nil.
func (s *ServiceInfo) UpdateServerStatus(server, status string) {
	s.serverStatusMu.Lock()
	defer s.serverStatusMu.Unlock()

	if s.serverStatus == nil {
		s.serverStatus = make(map[string]string)
	}
	s.serverStatus[server] = status
}

// GetAllStatus returns all the statuses of all the servers in ServiceInfo.
// It is the responsibility of the caller to check that s is not nil
func (s *ServiceInfo) GetAllStatus() map[string]string {
	s.serverStatusMu.RLock()
	defer s.serverStatusMu.RUnlock()

	if len(s.serverStatus) == 0 {
		return nil
	}

	allStatus := make(map[string]string, len(s.serverStatus))
	for k, v := range s.serverStatus {
		allStatus[k] = v
	}
	return allStatus
}

// RouterInfo holds information about a currently running HTTP router.
type RouterInfo struct {
	*dynamic.Router // dynamic configuration
	// Err contains all the errors that occurred during router's creation.
	Err []string `json:"error,omitempty"`
	// Status reports whether the router is disabled, in a warning state, or all good (enabled).
	// If not in "enabled" state, the reason for it should be in the list of Err.
	// It is the caller's responsibility to set the initial status.
	Status string   `json:"status,omitempty"`
	Using  []string `json:"using,omitempty"` // Effective entry points used by that router.
}

// AddError adds err to r.Err, if it does not already exist.
// If critical is set, r is marked as disabled.
func (r *RouterInfo) AddError(err error, critical bool) {
	for _, value := range r.Err {
		if value == err.Error() {
			return
		}
	}

	r.Err = append(r.Err, err.Error())
	if critical {
		r.Status = StatusDisabled
		return
	}

	// only set it to "warning" if not already in a worse state
	if r.Status != StatusDisabled {
		r.Status = StatusWarning
	}
}
