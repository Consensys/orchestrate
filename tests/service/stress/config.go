package stress

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/stress/utils"
)

func init() {
	viper.SetDefault(IterationViperKey, iterationDefault)
	_ = viper.BindEnv(IterationViperKey, iterationsEnv)
	viper.SetDefault(ArtifactPathViperKey, artifactPathDefault)
	_ = viper.BindEnv(ArtifactPathViperKey, artifactPathEnv)
	viper.SetDefault(ConcurrencyViperKey, concurrencyDefault)
	_ = viper.BindEnv(ConcurrencyViperKey, concurrencyEnv)
	viper.SetDefault(globalDataViperKey, globalDataDefault)
	_ = viper.BindEnv(globalDataViperKey, globalDataEnv)
}

// InitFlags register Cucumber flags
func InitFlags(f *pflag.FlagSet) {
	Iterations(f)
	Concurrency(f)
	ArtifactPath(f)
	GlobalData(f)
}

const (
	iterationFlag     = "stress-iterations"
	IterationViperKey = "stress.iterations"
	iterationDefault  = 10
	iterationsEnv     = "STRESS_ITERATIONS"
)

// Randomize register flag for randomize feature tests
func Iterations(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Number of test iteration execute per stress test unit
Environment variable: %q`, iterationsEnv)
	f.Int(iterationFlag, iterationDefault, desc)
	_ = viper.BindPFlag(IterationViperKey, f.Lookup(iterationFlag))
}

const (
	concurrencyFlag     = "stress-concurrency"
	ConcurrencyViperKey = "stress.concurrency"
	concurrencyDefault  = 1
	concurrencyEnv      = "STRESS_CONCURRENCY"
)

// Randomize register flag for randomize feature tests
func Concurrency(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Number of parallel threads spawn to accomplish the iterations
Environment variable: %q`, concurrencyEnv)
	f.Int(concurrencyFlag, concurrencyDefault, desc)
	_ = viper.BindPFlag(ConcurrencyViperKey, f.Lookup(concurrencyFlag))
}

var (
	artifactPathFlag     = "artifacts-paths"
	ArtifactPathViperKey = "artifacts.paths"
	artifactPathDefault  = []string{"/artifacts"}
	artifactPathEnv      = "ARTIFACTS_PATH"
)

// Artifact paths register flag for Godog Paths Option
func ArtifactPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`All artifact files path
Environment variable: %q`, artifactPathEnv)
	f.StringSlice(artifactPathFlag, artifactPathDefault, desc)
	_ = viper.BindPFlag(ArtifactPathViperKey, f.Lookup(artifactPathFlag))
}

var (
	globalDataFlag     = "stress-data"
	globalDataViperKey = "stress.data"
	globalDataDefault  = "{}"
	globalDataEnv      = "TEST_GLOBAL_DATA"
)

// Aliases register flag for aliases
func GlobalData(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Environment test data required by test (e.g chain.primary:888)
Environment variable: %q`, globalDataEnv)
	f.String(globalDataFlag, globalDataDefault, desc)
	_ = viper.BindPFlag(globalDataViperKey, f.Lookup(globalDataFlag))
}

type Config struct {
	ArtifactPath string
	Iterations   int
	Concurrency  int
	gData        utils.GlobalData
}

func InitConfig(vipr *viper.Viper) (*Config, error) {
	gd := utils.GlobalData{}
	raw := vipr.GetString(globalDataViperKey)
	err := json.Unmarshal([]byte(raw), &gd)
	if err != nil {
		return nil, err
	}

	return &Config{
		ArtifactPath: vipr.GetString(ArtifactPathViperKey),
		Iterations:   vipr.GetInt(IterationViperKey),
		Concurrency:  vipr.GetInt(ConcurrencyViperKey),
		gData:        gd,
	}, nil
}
