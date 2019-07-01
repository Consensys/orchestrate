package chain

import (
	"github.com/blang/semver"
	errors "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
)

var (
	component        = "types.chain"
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

func (p *Protocol) Version() (*semver.Version, error) {
	ver, err := semver.Make(p.Tag)
	if err != nil {
		return nil, errors.InvalidFormatError("invalid semver %q (%v)", p.Tag, err).SetComponent(component)
	}
	return &ver, nil
}

// IsTessera indicates wether if protocol is Quorum and support Tessera
func (p *Protocol) IsTessera() (bool, error) {
	if !p.IsQuorum() {
		return false, nil
	}

	ver, err := p.Version()
	if err != nil {
		return false, err
	}

	return ver.GE(tesseraVersion), nil
}
