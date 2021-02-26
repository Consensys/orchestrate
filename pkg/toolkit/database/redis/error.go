package redis

import (
	"github.com/ConsenSys/orchestrate/pkg/errors"
	ierror "github.com/ConsenSys/orchestrate/pkg/types/error"
	"github.com/gomodule/redigo/redis"
)

func parseRedisError(err error, msg string) *ierror.Error {
	if err == nil {
		return nil
	}

	switch {
	case err == redis.ErrNil:
		return errors.NotFoundError(msg).AppendReason(err.Error())
	default:
		return errors.RedisConnectionError(msg).AppendReason(err.Error())
	}
}
