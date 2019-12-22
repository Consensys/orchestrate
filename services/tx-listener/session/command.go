package session

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type Type int

const (
	START Type = iota
	STOP
)

type Command struct {
	Type Type
	Node *dynamic.Node
}

func CompareConfiguration(oldConfig, newConfig *dynamic.Configuration) []*Command {
	return []*Command{}
}
