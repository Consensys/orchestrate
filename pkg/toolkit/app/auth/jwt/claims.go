package jwt

import (
	"encoding/json"
	"strings"

	"github.com/golang-jwt/jwt"
)

type OrchestrateClaims struct {
	TenantID string `json:"tenant_id"`
}

type Claims struct {
	jwt.MapClaims
	Orchestrate OrchestrateClaims

	// Configurable JWT claims namespace for Orchestrate
	namespace string
}

func (c *Claims) UnmarshalJSON(b []byte) error {
	// First Unmarshal JWT entries
	err := json.Unmarshal(b, &c.MapClaims)
	if err != nil {
		return err
	}

	// Second Unmarshal Orchestrate entries
	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(b, &objmap)
	if err != nil {
		return err
	}

	if raw, ok := objmap[c.namespace]; c.namespace != "" && ok {
		err = json.Unmarshal(*raw, &c.Orchestrate)
		if err != nil {
			return err
		}
	} else {
		_, c.Orchestrate.TenantID = extractUsernameAndTenant(c.MapClaims["sub"].(string))
	}

	return nil
}

func extractUsernameAndTenant(sub string) (username, tenant string) {
	if !strings.Contains(sub, usernameTenantSeparator) {
		return "", sub
	}

	pieces := strings.Split(sub, usernameTenantSeparator)
	return pieces[1], pieces[0]
}
