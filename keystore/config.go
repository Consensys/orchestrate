package keystore

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(secretPkeyViperKey, secretPkeyDefault)
	viper.BindEnv(secretPkeyViperKey, secretPkeyEnv)
}

// InitFlags initialize flags
func InitFlags(f *pflag.FlagSet) {
	SecretPkeys(f)
}

var (
	secretPkeyFlag     = "secret-pkey"
	secretPkeyViperKey = "secret.pkeys"
	secretPkeyDefault  = []string{
		"56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E",
		"5FBB50BFF6DFAD35C4A374C9237BA2F7EAED9C6868E0108CB259B62D68029B1A",
		"86B021CCB810F26A30445B85F71E4C1596A11A97DDF9B9E348AC93D1DA6735BC",
		"DD614C3B343E1B6DBD1B2811D4F146CC90337DEEF96AB97C353578E871B19D5E",
		"425D92F63A836F890F1690B34B6A25C2971EF8D035CD8EA8592FD1069BD151C6",
		"C4B172E72033581BC41C36FA0448FCF031E9A31C4A3E300E541802DFB7248307",
		"706CC0876DA4D52B6DCE6F5A0FF210AEFCD51DE9F9CFE7D1BF7B385C82A06B8C",
		"1476C66DE79A57E8AB4CADCECCBE858C99E5EDF3BFFEA5404B15322B5421E18C",
		"A2426FE76ECA2AA7852B95A2CE9CC5CC2BC6C05BB98FDA267F2849A7130CF50D",
		"41B9C5E497CFE6A1C641EFCA314FF84D22036D1480AF5EC54558A5EDD2FEAC03",
	}
	secretPkeyEnv = "SECRET_PKEY"
)

// SecretPkeys register flag for Vault accounts
func SecretPkeys(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Private keys to pre-register in key store
Environment variable: %q`, secretPkeyEnv)
	f.StringSlice(secretPkeyFlag, secretPkeyDefault, desc)
	viper.BindPFlag(secretPkeyViperKey, f.Lookup(secretPkeyFlag))
}
