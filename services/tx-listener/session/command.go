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

		if !isEqualNode(newConfig.Nodes[k], oldConfig.Nodes[k]) {
			command := &Command{
				Type: UPDATE,
				Node: newConfig.Nodes[k],
			}
			commands = append(commands, command)
		}
	}

	return commands
}

func isEqualNode(node1, node2 *dynamic.Node) bool {
	return node1.TenantID == node2.TenantID &&
		node1.Name == node2.Name &&
		node1.URL == node2.URL &&
		node1.Listener.Depth == node2.Listener.Depth &&
		node1.Listener.Backoff == node2.Listener.Backoff
}
