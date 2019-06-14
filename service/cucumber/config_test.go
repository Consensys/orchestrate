package cucumber

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestShowStepDefinitions(t *testing.T) {
	name := "cucumber.showstepdefinitions"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	ShowStepDefinitions(flgs)

	// Test default
	expected := cucumberShowStepDefinitionsDefault
	assert.Equal(t, expected, viper.GetBool(name), "Default config should match")

	// Test environment variable
	os.Setenv("CUCUMBER_SHOWSTEPDEFINITION", "1")

	expected = true
	assert.Equal(t, expected, viper.GetBool(name), "Changing env var should change ShowStepDefinitions")
	os.Unsetenv("CUCUMBER_SHOWSTEPDEFINITION")

	// Test flags
	args := []string{
		"--cucumber-showstepdefinitions",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = true
	assert.Equal(t, expected, viper.GetBool(name), "Changing flags should change ShowStepDefinitions")
}

func TestRandomize(t *testing.T) {
	name := "cucumber.randomize"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Randomize(flgs)

	// Test default
	expected := cucumberRandomizeDefault
	assert.Equal(t, expected, viper.GetInt(name), "Default config should match")

	// Test environment variable
	os.Setenv("CUCUMBER_RANDOMIZE", "-1")
	expected = -1
	assert.Equal(t, expected, viper.GetInt(name), "Changing env var should change Randomize")
	os.Unsetenv("CUCUMBER_RANDOMIZE")

	// Test flags
	args := []string{
		"--cucumber-randomize=10",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = 10
	assert.Equal(t, expected, viper.GetInt(name), "Changing flags should change Randomize")
}

func TestStopOnFailure(t *testing.T) {
	name := "cucumber.stoponfailure"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	StopOnFailure(flgs)

	// Test default
	expected := cucumberStopOnFailureDefault
	assert.Equal(t, expected, viper.GetBool(name), "Default config should match")

	// Test environment variable
	os.Setenv("CUCUMBER_STOPONFAILURE", "1")
	expected = true
	assert.Equal(t, expected, viper.GetBool(name), "Changing env var should change StopOnFailure")
	os.Unsetenv("CUCUMBER_STOPONFAILURE")

	// Test flags
	args := []string{
		"--cucumber-stoponfailure",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = true
	assert.Equal(t, expected, viper.GetBool(name), "Changing flags should change StopOnFailure")
}

func TestStrict(t *testing.T) {
	name := "cucumber.strict"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Strict(flgs)

	// Test default
	expected := cucumberStrictDefault
	assert.Equal(t, expected, viper.GetBool(name), "Default config should match")

	// Test environment variable
	os.Setenv("CUCUMBER_STRICT", "1")
	expected = true
	assert.Equal(t, expected, viper.GetBool(name), "Changing env var should change Strict")
	os.Unsetenv("CUCUMBER_STRICT")

	// Test flags
	args := []string{
		"--cucumber-strict",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = true
	assert.Equal(t, expected, viper.GetBool(name), "Changing flags should change Strict")
}

func TestNoColors(t *testing.T) {
	name := "cucumber.nocolors"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	NoColors(flgs)

	// Test default
	expected := cucumberNoColorsDefault
	assert.Equal(t, expected, viper.GetBool(name), "Default config should match")

	// Test environment variable
	os.Setenv("CUCUMBER_NOCOLORS", "1")
	expected = true
	assert.Equal(t, expected, viper.GetBool(name), "Changing env var should change NoColors")
	os.Unsetenv("CUCUMBER_NOCOLORS")

	// Test flags
	args := []string{
		"--cucumber-nocolors",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = true
	assert.Equal(t, expected, viper.GetBool(name), "Changing flags should change NoColors")
}

func TestTags(t *testing.T) {
	name := "cucumber.tags"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Tags(flgs)

	// Test default
	expected := cucumberTagsDefault
	assert.Equal(t, expected, viper.GetString(name), "Default config should match")

	// Test environment variable
	os.Setenv("CUCUMBER_TAGS", "test")
	expected = "test"
	assert.Equal(t, expected, viper.GetString(name), "Changing env var should change Tags")
	os.Unsetenv("CUCUMBER_TAGS")

	// Test flags
	args := []string{
		"--cucumber-tags=test",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)
	expected = "test"
	assert.Equal(t, expected, viper.GetString(name), "Changing flags should change Tags")
}

func TestFormat(t *testing.T) {
	name := "cucumber.format"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Format(flgs)

	// Test default
	expected := cucumberFormatDefault
	assert.Equal(t, expected, viper.GetString(name), "Default config should match")

	// Test environment variable
	os.Setenv("CUCUMBER_FORMAT", "test")
	expected = "test"
	assert.Equal(t, expected, viper.GetString(name), "Changing env var should change Format")
	os.Unsetenv("CUCUMBER_FORMAT")

	// Test flags
	args := []string{
		"--cucumber-format=test",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)
	expected = "test"
	assert.Equal(t, expected, viper.GetString(name), "Changing flags should change Format")
}

func TestConcurrency(t *testing.T) {
	name := "cucumber.concurrency"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Concurrency(flgs)

	// Test default
	expected := cucumberConcurrencyDefault
	assert.Equal(t, expected, viper.GetInt(name), "Default config should match")

	// Test environment variable
	os.Setenv("CUCUMBER_CONCURRENCY", "20")
	expected = 20
	assert.Equal(t, expected, viper.GetInt(name), "Changing env var should change Concurrency")
	os.Unsetenv("CUCUMBER_CONCURRENCY")

	// Test flags
	args := []string{
		"--cucumber-concurrency=10",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = 10
	assert.Equal(t, expected, viper.GetInt(name), "Changing flags should change Concurrency")
}

func TestPaths(t *testing.T) {
	name := "cucumber.paths"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Paths(flgs)

	// Test default
	expected := cucumberPathsDefault
	assert.Equal(t, expected, viper.GetStringSlice(name), "Default config should match")

	// Test environment variable
	os.Setenv("CUCUMBER_PATHS", "test1 test2")
	expected = []string{
		"test1",
		"test2",
	}
	assert.Equal(t, expected, viper.GetStringSlice(name), "Changing env var should change Paths")
	os.Unsetenv("CUCUMBER_PATHS")

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
	assert.Equal(t, expected, viper.GetStringSlice(name), "Changing flags should change Paths")
}

func TestOutputPath(t *testing.T) {
	name := "cucumber.outputpath"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	OutputPath(flgs)

	// Test default
	expected := cucumberOutputPathDefault
	assert.Equal(t, expected, viper.GetString(name), "Default config should match")

	// Test environment variable
	os.Setenv("CUCUMBER_OUTPUTPATH", "test")
	expected = "test"
	assert.Equal(t, expected, viper.GetString(name), "Changing env var should change OutputPath")
	os.Unsetenv("CUCUMBER_OUTPUTPATH")

	// Test flags
	args := []string{
		"--cucumber-outputpath=test",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)
	expected = "test"
	assert.Equal(t, expected, viper.GetString(name), "Changing flags should change OutputPath")
}

func TestInitFlags(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	InitFlags(flgs)
}
