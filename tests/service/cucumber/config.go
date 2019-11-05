package cucumber

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(ShowStepDefinitionsViperKey, showStepDefinitionsDefault)
	_ = viper.BindEnv(ShowStepDefinitionsViperKey, showStepDefinitionsEnv)
	viper.SetDefault(RandomizeViperKey, randomizeDefault)
	_ = viper.BindEnv(RandomizeViperKey, randomizeEnv)
	viper.SetDefault(StopOnFailureViperKey, stopOnFailureDefault)
	_ = viper.BindEnv(StopOnFailureViperKey, stopOnFailureEnv)
	viper.SetDefault(StrictViperKey, strictDefault)
	_ = viper.BindEnv(StrictViperKey, strictEnv)
	viper.SetDefault(NoColorsViperKey, noColorsDefault)
	_ = viper.BindEnv(NoColorsViperKey, noColorsEnv)
	viper.SetDefault(TagsViperKey, tagsDefault)
	_ = viper.BindEnv(TagsViperKey, tagsEnv)
	viper.SetDefault(FormatViperKey, formatDefault)
	_ = viper.BindEnv(FormatViperKey, formatEnv)
	viper.SetDefault(ConcurrencyViperKey, concurrencyDefault)
	_ = viper.BindEnv(ConcurrencyViperKey, concurrencyEnv)
	viper.SetDefault(PathsViperKey, pathsDefault)
	_ = viper.BindEnv(PathsViperKey, pathsEnv)
	viper.SetDefault(OutputPathViperKey, outputPathDefault)
	_ = viper.BindEnv(OutputPathViperKey, outputPathEnv)
}

// InitFlags register Cucumber flags
func InitFlags(f *pflag.FlagSet) {
	ShowStepDefinitions(f)
	Randomize(f)
	StopOnFailure(f)
	Strict(f)
	NoColors(f)
	Tags(f)
	Format(f)
	Concurrency(f)
	Paths(f)
	OutputPath(f)
}

const (
	showStepDefinitionsFlag     = "cucumber-showstepdefinitions"
	ShowStepDefinitionsViperKey = "cucumber.showstepdefinitions"
	showStepDefinitionsDefault  = false
	showStepDefinitionsEnv      = "CUCUMBER_SHOWSTEPDEFINITION"
)

// ShowStepDefinitions register flag for Godog ShowStepDefinitions Option
func ShowStepDefinitions(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Print step definitions found and exit
Environment variable: %q`, showStepDefinitionsEnv)
	f.Bool(showStepDefinitionsFlag, showStepDefinitionsDefault, desc)
	_ = viper.BindPFlag(ShowStepDefinitionsViperKey, f.Lookup(showStepDefinitionsFlag))
}

const (
	randomizeFlag     = "cucumber-randomize"
	RandomizeViperKey = "cucumber.randomize"
	randomizeDefault  = -1
	randomizeEnv      = "CUCUMBER_RANDOMIZE"
)

// Randomize register flag for randomize feature tests
func Randomize(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Seed to randomize feature tests. The default value of -1 means to have a random seed. 0 means do not randomize
Environment variable: %q`, randomizeEnv)
	f.Int(randomizeFlag, randomizeDefault, desc)
	_ = viper.BindPFlag(RandomizeViperKey, f.Lookup(randomizeFlag))
}

const (
	stopOnFailureFlag     = "cucumber-stoponfailure"
	StopOnFailureViperKey = "cucumber.stoponfailure"
	stopOnFailureDefault  = false
	stopOnFailureEnv      = "CUCUMBER_STOPONFAILURE"
)

// StopOnFailure register flag for Godog StopOnFailure Option
func StopOnFailure(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Stops on the first failure
Environment variable: %q`, stopOnFailureEnv)
	f.Bool(stopOnFailureFlag, stopOnFailureDefault, desc)
	_ = viper.BindPFlag(StopOnFailureViperKey, f.Lookup(stopOnFailureFlag))
}

const (
	strictFlag     = "cucumber-strict"
	StrictViperKey = "cucumber.strict"
	strictDefault  = false
	strictEnv      = "CUCUMBER_STRICT"
)

// Strict register flag for Godog Strict Option
func Strict(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Fail suite when there are pending or undefined steps
Environment variable: %q`, strictEnv)
	f.Bool(strictFlag, strictDefault, desc)
	_ = viper.BindPFlag(StrictViperKey, f.Lookup(strictFlag))
}

const (
	noColorsFlag     = "cucumber-nocolors"
	NoColorsViperKey = "cucumber.nocolors"
	noColorsDefault  = false
	noColorsEnv      = "CUCUMBER_NOCOLORS"
)

// NoColors register flag for Godog NoColors Option
func NoColors(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Forces ansi color stripping
Environment variable: %q`, noColorsEnv)
	f.Bool(noColorsFlag, noColorsDefault, desc)
	_ = viper.BindPFlag(NoColorsViperKey, f.Lookup(noColorsFlag))
}

const (
	tagsFlag     = "cucumber-tags"
	TagsViperKey = "cucumber.tags"
	tagsDefault  = ""
	tagsEnv      = "CUCUMBER_TAGS"
)

// Tags register flag for Godog Tags Option
func Tags(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Various filters for scenarios parsed from feature files
Environment variable: %q`, tagsEnv)
	f.String(tagsFlag, tagsDefault, desc)
	_ = viper.BindPFlag(TagsViperKey, f.Lookup(tagsFlag))
}

const (
	formatFlag     = "cucumber-format"
	FormatViperKey = "cucumber.format"
	formatDefault  = "pretty"
	formatEnv      = "CUCUMBER_FORMAT"
)

// Format register flag for Godog Format Option
func Format(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`The formatter name (events|junit|pretty|cucumber)
Environment variable: %q`, formatEnv)
	f.String(formatFlag, formatDefault, desc)
	_ = viper.BindPFlag(FormatViperKey, f.Lookup(formatFlag))
}

const (
	concurrencyFlag     = "cucumber-concurrency"
	ConcurrencyViperKey = "cucumber.concurrency"
	concurrencyDefault  = 1
	concurrencyEnv      = "CUCUMBER_CONCURRENCY"
)

// Concurrency register flag for Godog Concurrency Option
func Concurrency(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Concurrency rate, not all formatters accepts this
Environment variable: %q`, concurrencyEnv)
	f.Int(concurrencyFlag, concurrencyDefault, desc)
	_ = viper.BindPFlag(ConcurrencyViperKey, f.Lookup(concurrencyFlag))
}

var (
	pathsFlag     = "cucumber-paths"
	PathsViperKey = "cucumber.paths"
	pathsDefault  = []string{"features"}
	pathsEnv      = "CUCUMBER_PATHS"
)

// Paths register flag for Godog Paths Option
func Paths(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`All feature file paths
Environment variable: %q`, pathsEnv)
	f.StringSlice(pathsFlag, pathsDefault, desc)
	_ = viper.BindPFlag(PathsViperKey, f.Lookup(pathsFlag))
}

const (
	outputPathFlag     = "cucumber-outputpath"
	OutputPathViperKey = "cucumber.outputpath"
	outputPathDefault  = ""
	outputPathEnv      = "CUCUMBER_OUTPUTPATH"
)

// OutputPath register flag for Godog OutputPath Option
func OutputPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Where it should print the cucumber output (only works with cucumber format)
Environment variable: %q`, outputPathEnv)
	f.String(outputPathFlag, outputPathDefault, desc)
	_ = viper.BindPFlag(OutputPathViperKey, f.Lookup(outputPathFlag))
}

var (
	aliasesFlag     = "cucumber-aliases"
	AliasesViperKey = "cucumber.aliases"
	aliasesDefault  []string
	aliasesEnv      = "CUCUMBER_ALIAS"
)

// Aliases register flag for aliases
func Aliases(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Aliases for cucumber test scenarios (e.g chain.primary:888)
Environment variable: %q`, aliasesEnv)
	f.StringSlice(aliasesFlag, aliasesDefault, desc)
	_ = viper.BindPFlag(AliasesViperKey, f.Lookup(aliasesFlag))
}
