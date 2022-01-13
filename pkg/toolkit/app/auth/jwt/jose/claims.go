package jose

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/consensys/orchestrate/pkg/types/entities"
)

type CustomClaims struct {
	UserClaims    *entities.UserClaims
	userClaimPath string
}

func NewCustomClaims(path string) *CustomClaims {
	return &CustomClaims{
		userClaimPath: path,
	}
}

func (c *CustomClaims) UnmarshalJSON(data []byte) error {
	var res map[string]interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	if _, ok := res[c.userClaimPath]; ok {
		c.UserClaims = &entities.UserClaims{}
		bClaims, _ := json.Marshal(res[c.userClaimPath])
		if err := json.Unmarshal(bClaims, &c.UserClaims); err != nil {
			return errors.New("invalid custom claims format")
		}
	} else {
		return errors.New("missing custom claims data")
	}

	return nil
}

func (c *CustomClaims) Validate(_ context.Context) error {
	return nil
}
