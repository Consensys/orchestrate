package multitenancy

import (
	"fmt"
	"strings"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/types/entities"
)

const AuthMethodJWT = "JWT"
const AuthMethodAPIKey = "API-KEY"

type UserInfo struct {
	// AuthMode records the mode that succeeded to Authenticate the request ('tls', 'api-key', 'oidc' or '')
	AuthMode string

	AuthValue string

	// Tenant belonged by the user
	TenantID string

	// Subject identifies the user
	Username string

	// Tenants user has access to
	AllowedTenants []string

	Wildcard bool
}

func NewJWTUserInfo(claims *entities.UserClaims, token string) *UserInfo {
	user := DefaultUser()
	user.TenantID = claims.TenantID
	user.Username = claims.Username
	user.AuthMode = AuthMethodJWT
	user.AuthValue = token
	user.AllowedTenants = genAllowedTenants(claims.TenantID)

	if user.TenantID == WildcardTenant {
		user.TenantID = DefaultTenant
		user.Wildcard = true
	}

	return user
}

func NewAPIKeyUserInfo(apiKey string) *UserInfo {
	return &UserInfo{
		AuthMode:       AuthMethodAPIKey,
		AuthValue:      apiKey,
		TenantID:       DefaultTenant,
		Wildcard:       true,
		AllowedTenants: []string{WildcardTenant},
	}
}

func NewUserInfo(tenantID, username string) *UserInfo {
	return &UserInfo{
		Username:       username,
		TenantID:       tenantID,
		AllowedTenants: genAllowedTenants(tenantID),
	}
}

func NewInternalAdminUser() *UserInfo {
	return &UserInfo{
		AuthMode:       AuthMethodAPIKey,
		TenantID:       WildcardTenant,
		Wildcard:       true,
		Username:       WildcardOwner,
		AllowedTenants: []string{WildcardTenant},
	}
}

func DefaultUser() *UserInfo {
	return &UserInfo{
		Username:       "",
		TenantID:       DefaultTenant,
		AllowedTenants: []string{DefaultTenant},
	}
}

func (u *UserInfo) ImpersonateTenant(tenantID string) error {
	switch {
	case tenantID == "":
		return nil
	case tenantID == WildcardTenant:
		return errors.PermissionDeniedError("wildcard user cannot be impersonate")
	case tenantID == DefaultTenant:
		u.TenantID = DefaultTenant
		u.AllowedTenants = []string{DefaultTenant}
	case tenantID != u.TenantID && u.HasTenantAccess(tenantID):
		u.TenantID = tenantID
		u.AllowedTenants = genAllowedTenants(tenantID)
	case tenantID != u.TenantID:
		return errors.PermissionDeniedError("access to tenant %q forbidden", tenantID)
	}

	// In case tenant impersonation, we clean up sourced username
	u.Username = ""
	return nil
}

// Feature restricted only to internal http request via API-KEY
func (u *UserInfo) ImpersonateUsername(username string) error {
	if username == "" {
		return nil
	}

	if u.AuthMode == AuthMethodAPIKey {
		u.Username = username
		return nil
	}

	return errors.PermissionDeniedError("access to username %q forbidden", username)
}

func (u *UserInfo) HasTenantAccess(tenantID string) bool {
	if u.TenantID == tenantID {
		return true
	}

	for _, allowedTenantID := range u.AllowedTenants {
		if allowedTenantID == tenantID || allowedTenantID == WildcardTenant {
			return true
		}
	}

	return false
}

func (u *UserInfo) HasUsernameAccess(username string) bool {
	if username == "" || u.Username == WildcardOwner {
		return true
	}

	if u.Username == username {
		return true
	}

	return false
}

func genAllowedTenants(tenantID string) []string {
	if tenantID == WildcardTenant {
		return []string{WildcardTenant}
	}

	allowedTenants := []string{DefaultTenant}

	subTenants := strings.Split(tenantID, utils.AuthSeparator)
	if len(subTenants) > 1 {
		var tenantGroup = subTenants[0]
		allowedTenants = append(allowedTenants, tenantGroup)
		for _, subTenant := range subTenants[1:] {
			tenantGroup = fmt.Sprintf("%s%s%s", tenantGroup, utils.AuthSeparator, subTenant)
			allowedTenants = append(allowedTenants, tenantGroup)
		}
	} else {
		allowedTenants = append(allowedTenants, tenantID)
	}

	return allowedTenants
}
