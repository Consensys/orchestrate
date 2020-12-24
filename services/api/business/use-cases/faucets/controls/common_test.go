// +build unit

package controls

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/chainregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
	"sort"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
)

var (
	chains = []string{
		uuid.Must(uuid.NewV4()).String(),
		uuid.Must(uuid.NewV4()).String(),
		uuid.Must(uuid.NewV4()).String(),
	}
	chainURLs = []string{
		"http://chain1.url",
		"http://chain2.url",
		"http://chain3.url",
	}
	addresses = []string{
		ethcommon.HexToAddress("0xab").Hex(),
		ethcommon.HexToAddress("0xcd").Hex(),
		ethcommon.HexToAddress("0xef").Hex(),
	}
)

func newFaucetReq(candidates map[string]*entities.Faucet, chainUUID, chainURL, beneficiary string) *chainregistry.Request {
	return &chainregistry.Request{
		Chain: &models.Chain{
			UUID: chainUUID,
			URLs: []string{chainURL},
		},
		Beneficiary: beneficiary,
		Candidates:  candidates,
	}
}

func electFirstFaucetCandidate(candidates map[string]*entities.Faucet) *entities.Faucet {
	keys := make([]string, 0, len(candidates))
	for k := range candidates {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return candidates[keys[0]]
}
