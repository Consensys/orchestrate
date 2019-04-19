package services

// OffsetMarker is an interface to mark that a context has been processed
// Typically marking kafka offsets of message that have been consumed
type OffsetMarker interface {
	Mark(msg interface{}) error
}

// Producer produces object
type Producer interface {
	Produce(pb interface{}) error
}
