package entities

import (
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

type Account struct {
	Alias               string
	Address             ethcommon.Address
	PublicKey           string
	CompressedPublicKey string
	TenantID            string
	OwnerID             string
	Attributes          map[string]string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
