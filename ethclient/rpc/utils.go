package rpc

import (
	"strings"

	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
)

func ClientTypeParser(clientVersion string) types.ClientType {
	if strings.Contains(clientVersion, "pantheon") {
		return types.PantheonClient
	}

	// Quorum client version also contains "Geth" string
	if strings.Contains(clientVersion, "quorum") {
		return types.QuorumClient
	}

	return types.UnknownClient
}
