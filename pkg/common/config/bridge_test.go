package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestBridgeLinks(t *testing.T) {
	name := "bridge.links"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BridgeLinks(flgs)

	// Test default
	expected := []string{}
	if len(expected) != len(viper.GetStringSlice(name)) {
		t.Errorf("BridgeLinks #1: expected %v but got %v", expected, viper.GetStringSlice(name))
	}

	// Test environment variable
	os.Setenv("BRIDGE_LINKS", "addr0@chainID0<>addr0@chainID0 addr1@chainID1<>addr1@chainID1")
	expected = []string{
		"addr0@chainID0<>addr0@chainID0",
		"addr1@chainID1<>addr1@chainID1",
	}
	if len(expected) != len(viper.GetStringSlice(name)) {
		t.Errorf("BridgeLinks #2: expect %v but got %v", expected, viper.GetStringSlice(name))
	} else {
		for i, url := range viper.GetStringSlice(name) {
			if url != expected[i] {
				t.Errorf("BridgeLinks #2: expect %v but got %v", expected, viper.GetStringSlice(name))
			}
		}
	}

	// Test flags
	args := []string{
		"--bridge-links=addr0@chainID0<>addr0@chainID0",
		"--bridge-links=addr1@chainID1<>add1@chainID1,addr2@chainID2<>addr2@chainID2",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = []string{
		"addr0@chainID0<>addr0@chainID0",
		"addr1@chainID1<>add1@chainID1",
		"addr2@chainID2<>addr2@chainID2",
	}
	if len(expected) != len(viper.GetStringSlice(name)) {
		t.Errorf("BridgeLinks #3: expect %v but got %v", expected, viper.GetStringSlice(name))
	} else {
		for i, url := range viper.GetStringSlice(name) {
			if url != expected[i] {
				t.Errorf("BridgeLinks #3: expect %v but got %v", expected, viper.GetStringSlice(name))
			}
		}
	}
}

func TestBridgeMethodSignature(t *testing.T) {
	name := "bridge.methodsignature"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BridgeMethodSignature(flgs)

	// Test default
	expected := "RelayMessage(bytes,address,address)"
	if expected != viper.GetString(name) {
		t.Errorf("BridgeMethodSignature #1: expected %v but got %v", expected, viper.GetString(name))
	}

	// Test environment variable
	expected = "TestMethod(address,uint256)"
	os.Setenv("BRIDGE_METHODSIGNATURE", expected)
	if expected != viper.GetString(name) {
		t.Errorf("BridgeMethodSignature #2: expect %v but got %v", expected, viper.GetString(name))
	}

	// Test flags
	args := []string{
		"--bridge-methodsignature=TestMethod(address,uint256)",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	if expected != viper.GetString(name) {
		t.Errorf("BridgeMethodSignature #3: expect %v but got %v", expected, viper.GetString(name))
	}
}

func TestBridgeAuthority(t *testing.T) {
	name := "bridge.authority"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BridgeAuthority(flgs)

	// Test default
	expected := ""
	if viper.GetString(name) != expected {
		t.Errorf("BridgeAuthority #1: expected %v but got %v", expected, viper.GetString(name))
	}

	// Test environment variable
	expected = "0xTestAddress"
	os.Setenv("BRIDGE_AUTHORITY", "0xTestAddress")
	if expected != viper.GetString(name) {
		t.Errorf("BridgeAuthority #2: expect %v but got %v", expected, viper.GetString(name))
	}

	// Test flags
	args := []string{
		"--bridge-authority=0xTestAddress",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	if expected != viper.GetString(name) {
		t.Errorf("BridgeAuthority #3: expect %v but got %v", expected, viper.GetString(name))
	}
}
