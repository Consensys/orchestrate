package types

type ClientType int

const (
	UnknownClient ClientType = iota
	QuorumClient
	PantheonClient
)
