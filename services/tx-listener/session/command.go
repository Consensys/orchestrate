package session

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type Type int

const (
	START Type = iota
	STOP
	UPDATE
)

type Command struct {
	Type Type
	Node *dynamic.Node
}

func CompareConfiguration(oldConfig, newConfig *dynamic.Configuration) []*Command {
	var commands []*Command
	for _, node := range newConfig.Nodes {
		commands = append(commands, &Command{Type: START, Node: node})
	}
	return commands
}
