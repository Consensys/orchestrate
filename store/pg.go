package store

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/spf13/viper"
)

// NewPGOptions creates new postgres options
// TODO: should be moved in pkg.git/common
func NewPGOptions() *pg.Options {
	return &pg.Options{
		Addr:     fmt.Sprintf("%v:%v", viper.GetString("db.host"), viper.GetString("db.port")),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		Database: viper.GetString("db.database"),
		PoolSize: viper.GetInt("db.poolsize"),
	}
}
