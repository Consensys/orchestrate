package envelope_store //nolint:golint,stylecheck

import (
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
)

func TestStatusInfoStatus(t *testing.T) {
	info := &StatusInfo{Status: Status_STORED}
	assert.False(t, info.HasBeenSent(), "#1: HasBeenSent")
	assert.False(t, info.IsPending(), "#1: IsPending")
	assert.False(t, info.IsMined(), "#1: IsMined")
	assert.False(t, info.IsError(), "#1: IsError")

	info = &StatusInfo{Status: Status_PENDING}
	assert.True(t, info.HasBeenSent(), "#2: HasBeenSent")
	assert.True(t, info.IsPending(), "#2: IsPending")
	assert.False(t, info.IsMined(), "#2: IsMined")
	assert.False(t, info.IsError(), "#2: IsError")

	info = &StatusInfo{Status: Status_MINED}
	assert.True(t, info.HasBeenSent(), "#3: HasBeenSent")
	assert.False(t, info.IsPending(), "#3: IsPending")
	assert.True(t, info.IsMined(), "#3: IsMined")
	assert.False(t, info.IsError(), "#3: IsError")

	info = &StatusInfo{Status: Status_ERROR}
	assert.True(t, info.HasBeenSent(), "#4: HasBeenSent")
	assert.False(t, info.IsPending(), "#4: IsPending")
	assert.False(t, info.IsMined(), "#4: IsMined")
	assert.True(t, info.IsError(), "#4: IsError")
}

func TestStatusInfoTimes(t *testing.T) {
	info := &StatusInfo{}
	assert.True(t, info.StoredAtTime().IsZero(), "Nil StoredAt should return null time")
	assert.True(t, info.SentAtTime().IsZero(), "Nil SentAt should return null time")
	assert.True(t, info.MinedAtTime().IsZero(), "Nil MinedAt should return null time")
	assert.True(t, info.ErrorAtTime().IsZero(), "Nil ErrorAt should return null time")

	info = &StatusInfo{
		StoredAt: ptypes.TimestampNow(),
	}
	assert.False(t, info.StoredAtTime().IsZero(), "Nil StoredAt should return null time")
}
