package session

import (
	"reflect"

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
	commands := make([]*Command, 0)

	for k, v := range newConfig.Nodes {
		if oldConfig.Nodes[k] == nil {
			command := &Command{
				Type: START,
				Node: v,
			}
			commands = append(commands, command)
		}
	}

	for k, v := range oldConfig.Nodes {
		if newConfig.Nodes[k] == nil {
			command := &Command{
				Type: STOP,
				Node: v,
			}
			commands = append(commands, command)
			continue
		}

		if !reflect.DeepEqual(newConfig.Nodes[k], oldConfig.Nodes[k]) {
			command := &Command{
				Type: UPDATE,
				Node: newConfig.Nodes[k],
			}
			commands = append(commands, command)
		}
	}

	return commands
}
