package utils

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	pduration "github.com/golang/protobuf/ptypes/duration"
	ptimestamp "github.com/golang/protobuf/ptypes/timestamp"
)

// TimestampToTime translates a protobuf timestamp into a time.Time
func PTimestampToTime(ptstamp *ptimestamp.Timestamp) time.Time {
	if ptstamp == nil {
		return time.Time{}
	}

	t, err := ptypes.Timestamp(ptstamp)
	if err != nil {
		panic(err)
	}
	return t
}

// TimeToTimestamp translates a time.Time into a protobuf timestamp
func TimeToPTimestamp(t time.Time) *ptimestamp.Timestamp {
	if t.IsZero() {
		return nil
	}

	ptstamp, err := ptypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}

	return ptstamp
}

func PDurationToDuration(pdur *pduration.Duration) time.Duration {
	if pdur == nil {
		return time.Duration(0)
	}

	d, err := ptypes.Duration(pdur)
	if err != nil {
		panic(err)
	}

	return d
}

func DurationToPDuration(d time.Duration) *pduration.Duration {
	if int64(d) == 0 {
		return nil
	}
	return ptypes.DurationProto(d)
}
