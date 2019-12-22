package ethereum

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session"
	hook "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/hooks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/offset"
)

type Session struct {
	Node *dynamic.Node

	hook    hook.Hook
	offsets offset.Manager
}

func (s *Session) Run(ctx context.Context) error {
	log.WithFields(log.Fields{
		"node.id":       s.Node.ID,
		"node.tenantId": s.Node.TenantID,
		"node.name":     s.Node.Name,
	}).Debugf("Start session")

	// TODO: implement session run
	<-ctx.Done()

	log.WithFields(log.Fields{
		"node.id":       s.Node.ID,
		"node.tenantId": s.Node.TenantID,
		"node.name":     s.Node.Name,
	}).Debugf("Session complete")
	return nil
}

type SessionBuilder struct {
	hook    hook.Hook
	offsets offset.Manager
}

func NewSessionBuilder(hk hook.Hook, offsets offset.Manager) *SessionBuilder {
	return &SessionBuilder{
		hook:    hk,
		offsets: offsets,
	}
}

func (b *SessionBuilder) NewSession(node *dynamic.Node) (session.Session, error) {
	return &Session{
		hook:    b.hook,
		offsets: b.offsets,
	}, nil
}
