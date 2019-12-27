package listenerv1

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type Provider struct {
	ec ethclient.ChainSyncReaderV2

	conf *Config
}

func NewProvider(conf *Config, ec ethclient.ChainSyncReaderV2) *Provider {
	return &Provider{
		conf: conf,
		ec:   ec,
	}
}

func (p *Provider) Run(ctx context.Context, configInput chan<- *dynamic.Message) error {
	msg := &dynamic.Message{
		Provider: "listener-v1",
		Configuration: &dynamic.Configuration{
			Nodes: make(map[string]*dynamic.Node),
		},
	}

	for _, url := range p.conf.URLs {
		// Create node
		node := &dynamic.Node{
			URL:    url,
			Active: true,
		}

		// Compute node ID from chain ID
		chainID, err := p.ec.Network(ctx, url)
		if err != nil {
			continue
		}
		node.ID = chainID.Text(10)

		// Set Listener configuration
		node.Listener = &dynamic.Listener{
			Depth:   p.conf.Depth,
			Backoff: p.conf.Backoff,
		}
		if pos, ok := p.conf.Start.Positions[chainID.Text(10)]; ok {
			node.Listener.BlockPosition = pos.BlockNumber
		} else {
			node.Listener.BlockPosition = p.conf.Start.Default.BlockNumber
		}
		msg.Configuration.Nodes[node.ID] = node
	}

	configInput <- msg

	return nil
}
