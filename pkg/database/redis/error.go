package redis

import (
	"github.com/gomodule/redigo/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/error"
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
