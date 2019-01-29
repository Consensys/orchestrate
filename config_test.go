package main

import (
	"testing"

	flags "github.com/jessevdk/go-flags"
)

func TestLoggerConfig(t *testing.T) {
	opts := LoggerConfig{}

	// Test default config
	args := []string{}
	args, err := flags.ParseArgs(&opts, args)

	if err != nil {
		t.Errorf("Unexpected error while parsing empty arg: %v", err)
	}

	expected := "text"
	if opts.Format != expected {
		t.Errorf("Default Format should be %q but got %q", expected, opts.Format)
	}

	expected = "debug"
	if opts.Level != expected {
		t.Errorf("Default Level should be %q but got %q", expected, opts.Format)
	}

	// Test args config
	args = []string{
		"--log-level=error",
		"--log-format=json",
	}

	args, err = flags.ParseArgs(&opts, args)
	if err != nil {
		t.Errorf("Unexpected error while parsing empty arg: %v", err)
	}

	expected = "json"
	if opts.Format != expected {
		t.Errorf("Format should be %q but got %q", expected, opts.Format)
	}

	expected = "error"
	if opts.Level != expected {
		t.Errorf("Level should be %q but got %q", expected, opts.Format)
	}
}

func TestWorkerConfig(t *testing.T) {
	opts := WorkerConfig{}

	// Test default config
	args := []string{}
	args, err := flags.ParseArgs(&opts, args)

	if err != nil {
		t.Errorf("Unexpected error while parsing empty arg: %v", err)
	}

	expected := uint(100)
	if opts.Slots != expected {
		t.Errorf("Default Slots should be %v but got %v", expected, opts.Slots)
	}

	// Test args config
	args = []string{
		"--worker-slots=1",
	}
	args, err = flags.ParseArgs(&opts, args)

	if err != nil {
		t.Errorf("Unexpected error while parsing empty arg: %v", err)
	}

	expected = uint(1)
	if opts.Slots != expected {
		t.Errorf("Default Slots should be %v but got %v", expected, opts.Slots)
	}
}

func TestEthClientConfig(t *testing.T) {
	opts := EthConfig{}

	// Test default config
	args := []string{}
	args, err := flags.ParseArgs(&opts, args)

	if err != nil {
		t.Errorf("Unexpected error while parsing empty arg: %v", err)
	}

	expected := []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
	}
	if len(expected) != len(opts.URLs) {
		t.Errorf("Default URLS should be %v but got %v", expected, opts.URLs)
	} else {
		for i, url := range opts.URLs {
			if url != expected[i] {
				t.Errorf("Default URLS should be %v but got %v", expected, opts.URLs)
			}
		}
	}

	// Test args config
	args = []string{
		"-e=http://localhost:8545",
		"--eth-client=http://localhost:7545",
	}
	args, err = flags.ParseArgs(&opts, args)

	if err != nil {
		t.Errorf("Unexpected error while parsing empty arg: %v", err)
	}

	expected = []string{
		"http://localhost:8545",
		"http://localhost:7545",
	}
	if len(expected) != len(opts.URLs) {
		t.Errorf("URLs should be %v but got %v", expected, opts.URLs)
	} else {
		for i, url := range opts.URLs {
			if url != expected[i] {
				t.Errorf("URLS should be %v but got %v", expected, opts.URLs)
			}
		}
	}
}
