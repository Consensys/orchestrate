package auth

import (
	"github.com/ConsenSys/orchestrate/pkg/auth/jwt"
	"github.com/ConsenSys/orchestrate/pkg/auth/key"
	"github.com/spf13/pflag"
)

func Flags(f *pflag.FlagSet) {
	key.Flags(f)
	jwt.Flags(f)
}
