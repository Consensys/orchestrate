package envelope_store //nolint:golint,stylecheck

import (
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
)

// HasBeenSent indicates whether a transaction has been sent
func (info *StatusInfo) HasBeenSent() bool {
	return info.GetStatus() != Status_STORED
}

// IsPending indicates whether transaction is pending
func (info *StatusInfo) IsPending() bool {
	return info.GetStatus() == Status_PENDING
}

// IsMined indicates whether transaction has been mined
func (info *StatusInfo) IsMined() bool {
	return info.GetStatus() == Status_MINED
}

// IsError indicates whether transaction is considered as an error
func (info *StatusInfo) IsError() bool {
	return info.GetStatus() == Status_ERROR
}

func (info *StatusInfo) StoredAtTime() time.Time {
	return utils.PTimestampToTime(info.GetStoredAt())
}

func (info *StatusInfo) SentAtTime() time.Time {
	return utils.PTimestampToTime(info.GetSentAt())
}

func (info *StatusInfo) MinedAtTime() time.Time {
	return utils.PTimestampToTime(info.GetMinedAt())
}

func (info *StatusInfo) ErrorAtTime() time.Time {
	return utils.PTimestampToTime(info.GetErrorAt())
}
