package auth

import (
	"github.com/spf13/pflag"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/jwt"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/key"
)

func Flags(f *pflag.FlagSet) {
	key.Flags(f)
	jwt.Flags(f)
}
