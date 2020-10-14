package identitymanager

import (
	"time"
)

type IdentityResponse struct {
	Alias               string            `json:"alias" example:"personal-account"`
	Address             string            `json:"address" example:"1abae27a0cbfb02945720425d3b80c7e09728534"`
	PublicKey           string            `json:"publicKey" example:"048e66b3e549818ea2cb354fb70749f6c8de8fa484f7530fc447d5fe80a1c424e4f5ae648d648c980ae7095d1efad87161d83886ca4b6c498ac22a93da5099014a"`
	CompressedPublicKey string            `json:"compressedPublicKey"`
	TenantID            string            `json:"tenantID" example:"foo"`
	Active              bool              `json:"active" example:"true"`
	Attributes          map[string]string `json:"attributes,omitempty"`
	CreatedAt           time.Time         `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt           time.Time         `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
}
