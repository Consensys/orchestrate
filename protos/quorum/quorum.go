package quorum

import (
	"github.com/blang/semver"
)

// TesseraVersion Quorum version from which Tessera is supported
var TesseraVersion = semver.MustParse("2.2.0")

// IsTessera indicates wether Tessera is supported by Quoru node we send the transaction to
func (q *Quorum) IsTessera() (bool, error) {
	ver, err := semver.Make(q.Version)
	if err != nil {
		return false, err
	}
	return ver.GE(TesseraVersion), nil
}
