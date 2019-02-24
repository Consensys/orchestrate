package infra

// import (
// 	"testing"

// 	flags "github.com/jessevdk/go-flags"
// )

// func TestFaucetConfig(t *testing.T) {
// 	opts := FaucetConfig{}

// 	// Test default config
// 	args := []string{}
// 	args, err := flags.ParseArgs(&opts, args)

// 	if err != nil {
// 		t.Errorf("Unexpected error while parsing empty arg: %v", err)
// 	}

// 	expectedAddresses := map[string]string{
// 		"3": "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
// 	}
// 	if len(expectedAddresses) != len(opts.Addresses) {
// 		t.Errorf("Default Addresses should be %v but got %v", expectedAddresses, opts.Addresses)
// 	} else {
// 		for k, v := range opts.Addresses {
// 			if v != expectedAddresses[k] {
// 				t.Errorf("Default Addresses should be %v but got %v", expectedAddresses, opts.Addresses)
// 			}
// 		}
// 	}

// 	expectedCD := "60s"
// 	if opts.CoolDownTime != expectedCD {
// 		t.Errorf("Default CoolDownTime should be %q but got %q", expectedCD, opts.CoolDownTime)
// 	}

// 	expectedBl := []string{
// 		"3-0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
// 	}
// 	if len(expectedBl) != len(opts.BlackList) {
// 		t.Errorf("Default BlackList should be %v but got %v", expectedBl, opts.BlackList)
// 	} else {
// 		for i, bl := range opts.BlackList {
// 			if bl != expectedBl[i] {
// 				t.Errorf("Default BlackList should be %v but got %v", expectedBl, opts.BlackList)
// 			}
// 		}
// 	}

// 	expectedMB := "200000000000000000"
// 	if opts.MaxBalance != expectedMB {
// 		t.Errorf("Default MaxBalance should be %q but got %q", expectedMB, opts.MaxBalance)
// 	}

// 	expectedTopic := "topic-tx-crafter"
// 	if opts.Topic != expectedTopic {
// 		t.Errorf("Default Topic should be %q but got %q", expectedTopic, opts.Topic)
// 	}

// 	// Test args config
// 	args = []string{
// 		"--faucet-address=2:0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
// 		"--faucet-address=5:0x664895b5fE3ddf049d2Fb508cfA03923859763C6",
// 		"--faucet-cd-time=1h45",
// 		"--faucet-blacklist=2-0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
// 		"--faucet-blacklist=5-0x664895b5fE3ddf049d2Fb508cfA03923859763C6",
// 		"--faucet-max-balance=100000",
// 		"--faucet-topic=topic-test",
// 	}

// 	args, err = flags.ParseArgs(&opts, args)
// 	if err != nil {
// 		t.Errorf("Unexpected error while parsing empty arg: %v", err)
// 	}

// 	expectedAddresses = map[string]string{
// 		"2": "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
// 		"5": "0x664895b5fE3ddf049d2Fb508cfA03923859763C6",
// 	}
// 	if len(expectedAddresses) != len(opts.Addresses) {
// 		t.Errorf("Addresses should be %v but got %v", expectedAddresses, opts.Addresses)
// 	} else {
// 		for k, v := range opts.Addresses {
// 			if v != expectedAddresses[k] {
// 				t.Errorf("Addresses should be %v but got %v", expectedAddresses, opts.Addresses)
// 			}
// 		}
// 	}

// 	expectedCD = "1h45"
// 	if opts.CoolDownTime != expectedCD {
// 		t.Errorf("CoolDownTime should be %q but got %q", expectedCD, opts.CoolDownTime)
// 	}

// 	expectedBl = []string{
// 		"2-0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
// 		"5-0x664895b5fE3ddf049d2Fb508cfA03923859763C6",
// 	}
// 	if len(expectedBl) != len(opts.BlackList) {
// 		t.Errorf("BlackList should be %v but got %v", expectedBl, opts.BlackList)
// 	} else {
// 		for i, bl := range opts.BlackList {
// 			if bl != expectedBl[i] {
// 				t.Errorf("BlackList should be %v but got %v", expectedBl, opts.BlackList)
// 			}
// 		}
// 	}

// 	expectedMB = "100000"
// 	if opts.MaxBalance != expectedMB {
// 		t.Errorf("MaxBalance should be %q but got %q", expectedMB, opts.MaxBalance)
// 	}

// 	expectedTopic = "topic-test"
// 	if opts.Topic != expectedTopic {
// 		t.Errorf("Topic should be %q but got %q", expectedTopic, opts.Topic)
// 	}
// }
