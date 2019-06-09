package cucumber

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(cucumberShowStepDefinitionsViperKey, cucumberShowStepDefinitionsDefault)
	_ = viper.BindEnv(cucumberShowStepDefinitionsViperKey, cucumberShowStepDefinitionsEnv)
	viper.SetDefault(cucumberRandomizeViperKey, cucumberRandomizeDefault)
	_ = viper.BindEnv(cucumberRandomizeViperKey, cucumberRandomizeEnv)
	viper.SetDefault(cucumberStopOnFailureViperKey, cucumberStopOnFailureDefault)
	_ = viper.BindEnv(cucumberStopOnFailureViperKey, cucumberStopOnFailureEnv)
	viper.SetDefault(cucumberStrictViperKey, cucumberStrictDefault)
	_ = viper.BindEnv(cucumberStrictViperKey, cucumberStrictEnv)
	viper.SetDefault(cucumberNoColorsViperKey, cucumberNoColorsDefault)
	_ = viper.BindEnv(cucumberNoColorsViperKey, cucumberNoColorsEnv)
	viper.SetDefault(cucumberTagsViperKey, cucumberTagsDefault)
	_ = viper.BindEnv(cucumberTagsViperKey, cucumberTagsEnv)
	viper.SetDefault(cucumberFormatViperKey, cucumberFormatDefault)
	_ = viper.BindEnv(cucumberFormatViperKey, cucumberFormatEnv)
	viper.SetDefault(cucumberConcurrencyViperKey, cucumberConcurrencyDefault)
	_ = viper.BindEnv(cucumberConcurrencyViperKey, cucumberConcurrencyEnv)
	viper.SetDefault(cucumberPathsViperKey, cucumberPathsDefault)
	_ = viper.BindEnv(cucumberPathsViperKey, cucumberPathsEnv)
	viper.SetDefault(cucumberOutputPathViperKey, cucumberOutputPathDefault)
	_ = viper.BindEnv(cucumberOutputPathViperKey, cucumberOutputPathEnv)
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

var (
	cucumberShowStepDefinitionsFlag     = "cucumber-showstepdefinitions"
	cucumberShowStepDefinitionsViperKey = "cucumber.showstepdefinitions"
	cucumberShowStepDefinitionsDefault  = false
	cucumberShowStepDefinitionsEnv      = "CUCUMBER_SHOWSTEPDEFINITION"
)

// ShowStepDefinitions register flag for Godog ShowStepDefinitions Option
func ShowStepDefinitions(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Print step definitions found and exit : %q`, cucumberShowStepDefinitionsEnv)
	f.Bool(cucumberShowStepDefinitionsFlag, cucumberShowStepDefinitionsDefault, desc)
	_ = viper.BindPFlag(cucumberShowStepDefinitionsViperKey, f.Lookup(cucumberShowStepDefinitionsFlag))
}

var (
	cucumberRandomizeFlag     = "cucumber-randomize"
	cucumberRandomizeViperKey = "cucumber.randomize"
	cucumberRandomizeDefault  = -1
	cucumberRandomizeEnv      = "CUCUMBER_RANDOMIZE"
)

// Randomize register flag for randomize feature tests
func Randomize(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Seed to randomize feature tests. The default value of -1 means to have a random seed. 0 means do not randomize : %q`, cucumberRandomizeEnv)
	f.Int(cucumberRandomizeFlag, cucumberRandomizeDefault, desc)
	_ = viper.BindPFlag(cucumberRandomizeViperKey, f.Lookup(cucumberRandomizeFlag))
}

var (
	cucumberStopOnFailureFlag     = "cucumber-stoponfailure"
	cucumberStopOnFailureViperKey = "cucumber.stoponfailure"
	cucumberStopOnFailureDefault  = false
	cucumberStopOnFailureEnv      = "CUCUMBER_STOPONFAILURE"
)

// StopOnFailure register flag for Godog StopOnFailure Option
func StopOnFailure(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Stops on the first failure : %q`, cucumberStopOnFailureEnv)
	f.Bool(cucumberStopOnFailureFlag, cucumberStopOnFailureDefault, desc)
	_ = viper.BindPFlag(cucumberStopOnFailureViperKey, f.Lookup(cucumberStopOnFailureFlag))
}

var (
	cucumberStrictFlag     = "cucumber-strict"
	cucumberStrictViperKey = "cucumber.strict"
	cucumberStrictDefault  = false
	cucumberStrictEnv      = "CUCUMBER_STRICT"
)

// Strict register flag for Godog Strict Option
func Strict(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Fail suite when there are pending or undefined steps : %q`, cucumberStrictEnv)
	f.Bool(cucumberStrictFlag, cucumberStrictDefault, desc)
	_ = viper.BindPFlag(cucumberStrictViperKey, f.Lookup(cucumberStrictFlag))
}

var (
	cucumberNoColorsFlag     = "cucumber-nocolors"
	cucumberNoColorsViperKey = "cucumber.nocolors"
	cucumberNoColorsDefault  = false
	cucumberNoColorsEnv      = "CUCUMBER_NOCOLORS"
)

// NoColors register flag for Godog NoColors Option
func NoColors(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Forces ansi color stripping : %q`, cucumberNoColorsEnv)
	f.Bool(cucumberNoColorsFlag, cucumberNoColorsDefault, desc)
	_ = viper.BindPFlag(cucumberNoColorsViperKey, f.Lookup(cucumberNoColorsFlag))
}

var (
	cucumberTagsFlag     = "cucumber-tags"
	cucumberTagsViperKey = "cucumber.tags"
	cucumberTagsDefault  = ""
	cucumberTagsEnv      = "CUCUMBER_TAGS"
)

// Tags register flag for Godog Tags Option
func Tags(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Various filters for scenarios parsed from feature files : %q`, cucumberTagsEnv)
	f.String(cucumberTagsFlag, cucumberTagsDefault, desc)
	_ = viper.BindPFlag(cucumberTagsViperKey, f.Lookup(cucumberTagsFlag))
}

var (
	cucumberFormatFlag     = "cucumber-format"
	cucumberFormatViperKey = "cucumber.format"
	cucumberFormatDefault  = "pretty"
	cucumberFormatEnv      = "CUCUMBER_FORMAT"
)

// Format register flag for Godog Format Option
func Format(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`The formatter name : %q`, cucumberFormatEnv)
	f.String(cucumberFormatFlag, cucumberFormatDefault, desc)
	_ = viper.BindPFlag(cucumberFormatViperKey, f.Lookup(cucumberFormatFlag))
}

var (
	cucumberConcurrencyFlag     = "cucumber-concurrency"
	cucumberConcurrencyViperKey = "cucumber.concurrency"
	cucumberConcurrencyDefault  = 1
	cucumberConcurrencyEnv      = "CUCUMBER_CONCURRENCY"
)

// Concurrency register flag for Godog Concurrency Option
func Concurrency(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Concurrency rate, not all formatters accepts this : %q`, cucumberConcurrencyEnv)
	f.Int(cucumberConcurrencyFlag, cucumberConcurrencyDefault, desc)
	_ = viper.BindPFlag(cucumberConcurrencyViperKey, f.Lookup(cucumberConcurrencyFlag))
}

var (
	cucumberPathsFlag     = "cucumber-paths"
	cucumberPathsViperKey = "cucumber.paths"
	cucumberPathsDefault  = []string{"features"}
	cucumberPathsEnv      = "CUCUMBER_PATHS"
)

// Paths register flag for Godog Paths Option
func Paths(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`All feature file paths : %q`, cucumberPathsEnv)
	f.StringSlice(cucumberPathsFlag, cucumberPathsDefault, desc)
	_ = viper.BindPFlag(cucumberPathsViperKey, f.Lookup(cucumberPathsFlag))
}

var (
	cucumberOutputPathFlag     = "cucumber-outputpath"
	cucumberOutputPathViperKey = "cucumber.outputpath"
	cucumberOutputPathDefault  = ""
	cucumberOutputPathEnv      = "CUCUMBER_OUTPUTPATH"
)

// OutputPath register flag for Godog OutputPath Option
func OutputPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Where it should print the cucumber output (only works with cucumber format): %q`, cucumberOutputPathEnv)
	f.String(cucumberOutputPathFlag, cucumberOutputPathDefault, desc)
	_ = viper.BindPFlag(cucumberOutputPathViperKey, f.Lookup(cucumberOutputPathFlag))
}
