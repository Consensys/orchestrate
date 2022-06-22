package api

import (
	"encoding/json"
	"time"

	"github.com/consensys/orchestrate/pkg/types/entities"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const DefaultAccountPageSize = 25

type AccountResponse struct {
	Alias               string            `json:"alias" example:"personal-account"`
	Address             ethcommon.Address `json:"address" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534" swaggertype:"string"`
	PublicKey           hexutil.Bytes     `json:"publicKey" example:"0x048e66b3e549818ea2cb354fb70749f6c8de8fa484f7530fc447d5fe80a1c424e4f5ae648d648c980ae7095d1efad87161d83886ca4b6c498ac22a93da5099014a" swaggertype:"string"`
	CompressedPublicKey hexutil.Bytes     `json:"compressedPublicKey" example:"0x048e66b3e549818ea2cb354fb70749f6c8de8fa484f7530fc447" swaggertype:"string"`
	TenantID            string            `json:"tenantID" example:"tenantFoo"`
	OwnerID             string            `json:"ownerID,omitempty" example:"foo"`
	StoreID             string            `json:"storeID,omitempty" example:"myQKMStoreID"`
	Attributes          map[string]string `json:"attributes,omitempty"`
	CreatedAt           time.Time         `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt           time.Time         `json:"updatedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
}

type AccountSearchResponse struct {
	Accounts []*AccountResponse `json:"accounts"`
	HasMore  bool               `json:"hasMore"`
}

type accountResponseJSON struct {
	Alias               string            `json:"alias"`
	Address             string            `json:"address"`
	PublicKey           string            `json:"publicKey"`
	CompressedPublicKey string            `json:"compressedPublicKey"`
	TenantID            string            `json:"tenantID"`
	OwnerID             string            `json:"ownerID,omitempty"`
	StoreID             string            `json:"storeID,omitempty"`
	Attributes          map[string]string `json:"attributes,omitempty"`
	CreatedAt           time.Time         `json:"createdAt"`
	UpdatedAt           time.Time         `json:"updatedAt,omitempty"`
}

func (a *AccountResponse) MarshalJSON() ([]byte, error) {
	res := &accountResponseJSON{
		Alias:               a.Alias,
		PublicKey:           a.PublicKey.String(),
		CompressedPublicKey: a.CompressedPublicKey.String(),
		Address:             a.Address.Hex(),
		TenantID:            a.TenantID,
		OwnerID:             a.OwnerID,
		StoreID:             a.StoreID,
		Attributes:          a.Attributes,
		CreatedAt:           a.CreatedAt,
		UpdatedAt:           a.UpdatedAt,
	}

	return json.Marshal(res)
}

func NewAccountResponse(acc *entities.Account) *AccountResponse {
	return &AccountResponse{
		Alias:               acc.Alias,
		Attributes:          acc.Attributes,
		Address:             acc.Address,
		PublicKey:           acc.PublicKey,
		CompressedPublicKey: acc.CompressedPublicKey,
		TenantID:            acc.TenantID,
		OwnerID:             acc.OwnerID,
		StoreID:             acc.StoreID,
		CreatedAt:           acc.CreatedAt,
		UpdatedAt:           acc.UpdatedAt,
	}
}
