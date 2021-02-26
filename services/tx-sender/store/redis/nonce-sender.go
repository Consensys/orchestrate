package redis

import (
	"github.com/ConsenSys/orchestrate/pkg/toolkit/database/redis"
	"github.com/ConsenSys/orchestrate/services/tx-sender/store"
)

type nonceSender struct {
	redis *redis.Client
}

// NewNonceSender creates a new mock NonceManager
func NewNonceSender(client *redis.Client) store.NonceSender {
	return &nonceSender{
		redis: client,
	}
}

const lastSentSuf = "last-sent"

func (ns *nonceSender) GetLastSent(key string) (nonce uint64, ok bool, err error) {
	return ns.redis.LoadUint64(computeKey(key, lastSentSuf))
}

func (ns *nonceSender) SetLastSent(key string, value uint64) error {
	return ns.redis.Set(computeKey(key, lastSentSuf), value)
}

func (ns *nonceSender) IncrLastSent(key string) error {
	return ns.redis.Incr(computeKey(key, lastSentSuf))
}

// IncrLastSent increment last sent nonce
func (ns *nonceSender) DeleteLastSent(key string) error {
	return ns.redis.Delete(computeKey(key, lastSentSuf))
}
