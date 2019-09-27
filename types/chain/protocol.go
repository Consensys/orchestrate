package chain

// IsPantheon indicates whether if protocol is Pantheon
func (p *Protocol) IsPantheon() bool {
	return p.GetType() == ProtocolType_PANTHEON_ORION
}

// IsConstellation indicates whether if protocol is Constellation with Quorum
func (p *Protocol) IsConstellation() bool {
	return p.GetType() == ProtocolType_QUORUM_CONSTELLATION
}

// IsTessera indicates whether if protocol is Tessera with Quorum
func (p *Protocol) IsTessera() bool {
	return p.GetType() == ProtocolType_QUORUM_TESSERA
}
