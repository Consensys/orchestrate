package engine

// Msg is an abstract interface supported by any kind of message handled by the engine
type Msg interface {
	// Entrypoint returns an indication on where the message comes from
	Entrypoint() string

	// Value returns value of the message
	Value() []byte

	// Key returns key associated to the message
	Key() []byte

	// Header return headers attached to the message
	Header() Header
}

// Header represents a key-value pairs in a message request.
type Header interface {
	// Add an entry to the header
	Add(key, value string)

	// Del deletes an header entry
	Del(key string)

	// Get retrieves an header entry
	Get(key string) string

	// Set an header entry
	Set(key, value string)
}
