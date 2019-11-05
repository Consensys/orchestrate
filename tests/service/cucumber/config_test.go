package cucumber

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestShowStepDefinitions(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	ShowStepDefinitions(flgs)

	// Test default
	expected := showStepDefinitionsDefault
	assert.Equal(t, expected, viper.GetBool(ShowStepDefinitionsViperKey), "Default config should match")

	// Test environment variable
	_ = os.Setenv("CUCUMBER_SHOWSTEPDEFINITION", "1")

	expected = true
	assert.Equal(t, expected, viper.GetBool(ShowStepDefinitionsViperKey), "Changing env var should change ShowStepDefinitions")
	_ = os.Unsetenv("CUCUMBER_SHOWSTEPDEFINITION")

	// Test flags
	args := []string{
		"--cucumber-showstepdefinitions",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = true
	assert.Equal(t, expected, viper.GetBool(ShowStepDefinitionsViperKey), "Changing flags should change ShowStepDefinitions")
}

func TestRandomize(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Randomize(flgs)

	// Test default
	expected := randomizeDefault
	assert.Equal(t, expected, viper.GetInt(RandomizeViperKey), "Default config should match")

	// Test environment variable
	_ = os.Setenv("CUCUMBER_RANDOMIZE", "-1")
	expected = -1
	assert.Equal(t, expected, viper.GetInt(RandomizeViperKey), "Changing env var should change Randomize")
	_ = os.Unsetenv("CUCUMBER_RANDOMIZE")

	// Test flags
	args := []string{
		"--cucumber-randomize=10",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = 10
	assert.Equal(t, expected, viper.GetInt(RandomizeViperKey), "Changing flags should change Randomize")
}

func TestStopOnFailure(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	StopOnFailure(flgs)

	// Test default
	expected := stopOnFailureDefault
	assert.Equal(t, expected, viper.GetBool(StopOnFailureViperKey), "Default config should match")

	// Test environment variable
	_ = os.Setenv("CUCUMBER_STOPONFAILURE", "1")
	expected = true
	assert.Equal(t, expected, viper.GetBool(StopOnFailureViperKey), "Changing env var should change StopOnFailure")
	_ = os.Unsetenv("CUCUMBER_STOPONFAILURE")

	// Test flags
	args := []string{
		"--cucumber-stoponfailure",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = true
	assert.Equal(t, expected, viper.GetBool(StopOnFailureViperKey), "Changing flags should change StopOnFailure")
}

func TestStrict(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Strict(flgs)

	// Test default
	expected := strictDefault
	assert.Equal(t, expected, viper.GetBool(StrictViperKey), "Default config should match")

	// Test environment variable
	_ = os.Setenv("CUCUMBER_STRICT", "1")
	expected = true
	assert.Equal(t, expected, viper.GetBool(StrictViperKey), "Changing env var should change Strict")
	_ = os.Unsetenv("CUCUMBER_STRICT")

	// Test flags
	args := []string{
		"--cucumber-strict",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = true
	assert.Equal(t, expected, viper.GetBool(StrictViperKey), "Changing flags should change Strict")
}

func TestNoColors(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	NoColors(flgs)

	// Test default
	expected := noColorsDefault
	assert.Equal(t, expected, viper.GetBool(NoColorsViperKey), "Default config should match")

	// Test environment variable
	_ = os.Setenv("CUCUMBER_NOCOLORS", "1")
	expected = true
	assert.Equal(t, expected, viper.GetBool(NoColorsViperKey), "Changing env var should change NoColors")
	_ = os.Unsetenv("CUCUMBER_NOCOLORS")

	// Test flags
	args := []string{
		"--cucumber-nocolors",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = true
	assert.Equal(t, expected, viper.GetBool(NoColorsViperKey), "Changing flags should change NoColors")
}

func TestTags(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Tags(flgs)

	// Test default
	expected := tagsDefault
	assert.Equal(t, expected, viper.GetString(TagsViperKey), "Default config should match")

	// Test environment variable
	_ = os.Setenv("CUCUMBER_TAGS", "test")
	expected = "test"
	assert.Equal(t, expected, viper.GetString(TagsViperKey), "Changing env var should change Tags")
	_ = os.Unsetenv("CUCUMBER_TAGS")

	// Test flags
	args := []string{
		"--cucumber-tags=test",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)
	expected = "test"
	assert.Equal(t, expected, viper.GetString(TagsViperKey), "Changing flags should change Tags")
}

func TestFormat(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Format(flgs)

	// Test default
	expected := formatDefault
	assert.Equal(t, expected, viper.GetString(FormatViperKey), "Default config should match")

	// Test environment variable
	_ = os.Setenv("CUCUMBER_FORMAT", "test")
	expected = "test"
	assert.Equal(t, expected, viper.GetString(FormatViperKey), "Changing env var should change Format")
	_ = os.Unsetenv("CUCUMBER_FORMAT")

	// Test flags
	args := []string{
		"--cucumber-format=test",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)
	expected = "test"
	assert.Equal(t, expected, viper.GetString(FormatViperKey), "Changing flags should change Format")
}

func TestConcurrency(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Concurrency(flgs)

	// Test default
	expected := concurrencyDefault
	assert.Equal(t, expected, viper.GetInt(ConcurrencyViperKey), "Default config should match")

	// Test environment variable
	_ = os.Setenv("CUCUMBER_CONCURRENCY", "20")
	expected = 20
	assert.Equal(t, expected, viper.GetInt(ConcurrencyViperKey), "Changing env var should change Concurrency")
	_ = os.Unsetenv("CUCUMBER_CONCURRENCY")

	// Test flags
	args := []string{
		"--cucumber-concurrency=10",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = 10
	assert.Equal(t, expected, viper.GetInt(ConcurrencyViperKey), "Changing flags should change Concurrency")
}

func TestPaths(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Paths(flgs)

	// Test default
	expected := pathsDefault
	assert.Equal(t, expected, viper.GetStringSlice(PathsViperKey), "Default config should match")

	// Test environment variable
	_ = os.Setenv("CUCUMBER_PATHS", "test1 test2")
	expected = []string{
		"test1",
		"test2",
	}
	assert.Equal(t, expected, viper.GetStringSlice(PathsViperKey), "Changing env var should change Paths")
	_ = os.Unsetenv("CUCUMBER_PATHS")

	// Test flags
	args := []string{
		"--cucumber-paths=test3",
		"--cucumber-paths=test4",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)
	expected = []string{
		"test3",
		"test4",
	}
	assert.Equal(t, expected, viper.GetStringSlice(PathsViperKey), "Changing flags should change Paths")
}

func TestOutputPath(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	OutputPath(flgs)

	// Test default
	expected := outputPathDefault
	assert.Equal(t, expected, viper.GetString(OutputPathViperKey), "Default config should match")

	// Test environment variable
	_ = os.Setenv("CUCUMBER_OUTPUTPATH", "test")
	expected = "test"
	assert.Equal(t, expected, viper.GetString(OutputPathViperKey), "Changing env var should change OutputPath")
	_ = os.Unsetenv("CUCUMBER_OUTPUTPATH")

	// Test flags
	args := []string{
		"--cucumber-outputpath=test",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)
	expected = "test"
	assert.Equal(t, expected, viper.GetString(OutputPathViperKey), "Changing flags should change OutputPath")
}

func TestInitFlags(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	InitFlags(flgs)
}
