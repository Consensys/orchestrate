package entities

// UserClaims represent raw claims extracted from an authentication method
type UserClaims struct {
	TenantID string `json:"tenant_id"`
	Username string `json:"username"`
}
