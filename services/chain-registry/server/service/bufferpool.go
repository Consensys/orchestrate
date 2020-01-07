// All scripts in server/service are highly inspired from Traefik server
// c.f. https://github.com/containous/traefik/tree/v2.0.5/pkg/server/service

package service

import "sync"

const bufferPoolSize = 32 * 1024

func newBufferPool() *bufferPool {
	return &bufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &buffer{b: make([]byte, bufferPoolSize)}
			},
		},
	}
}

type bufferPool struct {
	pool sync.Pool
}

type buffer struct {
	b []byte
}

func (b *bufferPool) Get() []byte {
	return b.pool.Get().(*buffer).b
}

func (b *bufferPool) Put(bytes []byte) {
	b.pool.Put(&buffer{b: bytes})
}
