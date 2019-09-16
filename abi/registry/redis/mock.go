package redis

import (
	"github.com/rafaeljusto/redigomock"
	remote "github.com/gomodule/redigo/redis"
)

// DialMock returns a mocked redis connexion
func DialMock(_, _ string, _ ...remote.DialOption) (remote.Conn, error) {
	return redigomock.NewConn(), nil
}