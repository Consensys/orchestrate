package chain

import (
	"github.com/blang/semver"
)

var (
	pantheonProtocol = "pantheon"
	quorumProtocol   = "quorum"
	tesseraVersion   = semver.MustParse("2.2.0")
)

// IsPantheon indicates wether if protocol is Pantheon
func (p *Protocol) IsPantheon() bool {
	return p.Name == pantheonProtocol
}

// IsQuorum indicates wether if protocol is Quorum
func (p *Protocol) IsQuorum() bool {
	return p.Name == quorumProtocol
}

// IsTessera indicates wether if protocol is Quorum and support Tessera
func (p *Protocol) IsTessera() (bool, error) {
	if !p.IsQuorum() {
		return false, nil
	}
	ver, err := semver.Make(p.Tag)
	if err != nil {
		return false, err
	}
	return ver.GE(tesseraVersion), nil
}
