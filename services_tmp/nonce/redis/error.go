package redis

import (
	"github.com/gomodule/redigo/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
)

func FromRedisError(err error) *ierror.Error {
	if err == nil {
		return nil
	}

	switch {
	case err == redis.ErrNil:
		return errors.NotFoundError(err.Error())
	default:
		return errors.InternalError(err.Error())
	}
}
