package auth

import (
	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/jwt/jose"
	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/key"
	"github.com/spf13/pflag"
)

func Flags(f *pflag.FlagSet) {
	key.Flags(f)
	jose.Flags(f)
}
