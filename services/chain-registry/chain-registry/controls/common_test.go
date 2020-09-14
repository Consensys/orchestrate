// +build unit

package controls

import (
	"sort"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/chainregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
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
	addresses = []ethcommon.Address{
		ethcommon.HexToAddress("0xab"),
		ethcommon.HexToAddress("0xcd"),
		ethcommon.HexToAddress("0xef"),
	}
)

func newFaucetReq(candidates map[string]chainregistry.Faucet, chainUUID, chainURL string, beneficiary ethcommon.Address) *chainregistry.Request {
	return &chainregistry.Request{
		Chain: &models.Chain{
			UUID: chainUUID,
			URLs: []string{chainURL},
		},
		Beneficiary: beneficiary,
		Candidates:  candidates,
	}
}

func electFirstFaucetCandidate(candidates map[string]chainregistry.Faucet) chainregistry.Faucet {
	keys := make([]string, 0, len(candidates))
	for k := range candidates {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return candidates[keys[0]]
}
