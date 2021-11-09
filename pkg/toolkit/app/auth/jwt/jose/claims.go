package jose

import "context"

type CustomClaims struct {
	TenantID string `json:"tenant_id"`
}

func (claims *CustomClaims) Validate(_ context.Context) error {
	// TODO: Apply validation on custom claims if needed, currently no validation is needed
	return nil
}
