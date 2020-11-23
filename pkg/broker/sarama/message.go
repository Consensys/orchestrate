package sarama

import (
	"bytes"

	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
)

// Msg is a wrapper struct for sarama.ConsumerMessage that implements engine.Msg
// As sarama.ConsumerMessage already have fields Key and Value that would conflict with
// the interface definition of engine.Msg, we resort in using a composite struct to
// solve the ambiguity
type Msg struct {
	sarama.ConsumerMessage
}

// Entrypoint returns the kafka topic of the message
func (m *Msg) Entrypoint() string { return m.Topic }

// Value returns the value, ie the data of the message
func (m *Msg) Value() []byte { return m.ConsumerMessage.Value }

// Key returns the key of the message
func (m *Msg) Key() []byte { return m.ConsumerMessage.Key }

// Header returns the header of the msg
func (m *Msg) Header() engine.Header { return (*Header)(&m.Headers) }

// Header is an alias for Header []*sarama.RecordHeader that implements engine.Header
type Header []*sarama.RecordHeader

// Get returns the value of an header entry by key. If the key does not exists
// it return "".
func (h Header) Get(key string) string {
	for _, rh := range h {
		if bytes.Equal([]byte(key), rh.Key) {
			return string(rh.Value)
		}
	}
	// If nothing was found, returns ""
	return ""
}

// Set adds a new key, value pair in the header or modify if it exists.
func (h *Header) Set(key, value string) {
	for index, rh := range *h {
		if bytes.Equal([]byte(key), rh.Key) {
			(*h)[index].Value = []byte(value)
			return
		}
	}
	// If the key does not exist in the header. Add a new one
	h.addUnsafe(key, value)
}

// Add a new key, value pair in the entry. If a identical key already exists.
// The function does nothing. This is the main difference with Set
func (h *Header) Add(key, value string) {
	for _, rh := range *h {
		if bytes.Equal([]byte(key), rh.Key) {
			return
		}
	}
	// If the key does not exist in the header. Add a new one
	h.addUnsafe(key, value)
}

// Adds a new key, value pair in the entry.
// Should not be called twice as it would create a duplicate entry.
// That's the reason why the method is private
func (h *Header) addUnsafe(key, value string) {
	*h = append(*h, &sarama.RecordHeader{
		Key:   []byte(key),
		Value: []byte(value),
	})
}

// Del removes an entry from the header.
// If duplicate keys exists, they are also removed.
func (h *Header) Del(key string) {
	for index, rh := range *h {
		if bytes.Equal([]byte(key), rh.Key) {
			l := len(*h)
			// Copy last element of array in delete element's index
			(*h)[index] = (*h)[l-1]
			// Remove the last element that we just copied
			*h = (*h)[:l-1]
			// Del
			return
		}
	}
	// If key does not exist, do nothing
}
