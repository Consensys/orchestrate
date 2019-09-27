package nonce

// Attributor allows to keep track of nonces as they are attributed to transactions
type Attributor interface {
	//  GetLastAttributed retrieves last attributed nonce
	GetLastAttributed(key string) (nonce uint64, ok bool, err error)

	// IncrLastAttributed Increment last attributed nonce
	IncrLastAttributed(key string) error

	// SetLastAttributed sets the last attributed nonce
	SetLastAttributed(key string, value uint64) error
}

// Sender allows to keep track of nonces as transactions are sent to blockchain
type Sender interface {
	// GetLastSent retrieves last sent nonce
	GetLastSent(key string) (nonce uint64, ok bool, err error)

	// IncrLastSent increment last sent nonce
	IncrLastSent(key string) error

	// SetLastSent sets last sent nonce
	SetLastSent(key string, value uint64) error

	// IsRecovering indicates whether we are recovering nonces
	IsRecovering(key string) (bool, error)

	// SetRevoring allows to set recovering status
	SetRecovering(key string, status bool) error
}

// Manager is an interface for NonceManagers
type Manager interface {
	Attributor
	Sender
}
