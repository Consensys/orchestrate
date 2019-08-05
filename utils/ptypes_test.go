package utils

import (
	"testing"
	"time"

	pduration "github.com/golang/protobuf/ptypes/duration"
	ptimestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
)

func TestPTimestampToTime(t *testing.T) {
	tstamp := PTimestampToTime(nil)
	t.Logf("%v", tstamp)
	assert.True(t, tstamp.Equal(time.Time{}), "Nil timestamp should transform into origin of time")

	tstamp = PTimestampToTime(&ptimestamp.Timestamp{})
	assert.True(t, tstamp.Equal(time.Unix(0, 0)), "TimeStamp should transform into origin of time")

	tstamp = PTimestampToTime(&ptimestamp.Timestamp{
		Nanos:   89754,
		Seconds: 100,
	})
	assert.Equal(t, 89754, tstamp.Nanosecond(), "Nanosecond should correct")
	assert.Equal(t, 40, tstamp.Second(), "Second should correct")
	assert.Equal(t, 1, tstamp.Minute(), "Minute should correct")
}

func TestPDurationToDuration(t *testing.T) {
	d := PDurationToDuration(&pduration.Duration{})
	assert.Equal(t, 0, int(d), "Nil should transform in 0")

	d = PDurationToDuration(&pduration.Duration{
		Nanos:   190,
		Seconds: 10,
	})
	assert.Equal(t, 10000000190, int(d), "Nil should transform in 0")
}
