package session

import (
	"github.com/consensys/orchestrate/services/tx-listener/dynamic"
)

type Type int

const (
	START Type = iota
	STOP
	UPDATE
)

type Command struct {
	Type  Type
	Chain *dynamic.Chain
}

func CompareConfiguration(oldConfig, newConfig *dynamic.Configuration) []*Command {
	commands := make([]*Command, 0)

	for k, v := range newConfig.Chains {
		if oldConfig.Chains[k] == nil {
			command := &Command{
				Type:  START,
				Chain: v,
			}
			commands = append(commands, command)
		}
	}

	for k, v := range oldConfig.Chains {
		if newConfig.Chains[k] == nil {
			command := &Command{
				Type:  STOP,
				Chain: v,
			}
			commands = append(commands, command)
			continue
		}

		if !isEqualChain(newConfig.Chains[k], oldConfig.Chains[k]) {
			command := &Command{
				Type:  UPDATE,
				Chain: newConfig.Chains[k],
			}
			commands = append(commands, command)
		}
	}

	return commands
}

func isEqualChain(chain1, chain2 *dynamic.Chain) bool {
	return chain1.TenantID == chain2.TenantID &&
		chain1.Name == chain2.Name &&
		chain1.URL == chain2.URL &&
		chain1.Listener.Depth == chain2.Listener.Depth &&
		chain1.Listener.Backoff == chain2.Listener.Backoff &&
		chain1.Listener.ExternalTxEnabled == chain2.Listener.ExternalTxEnabled
}
