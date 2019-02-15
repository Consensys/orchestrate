package main

import "testing"

func TestTranslateBlockNumber(t *testing.T) {
	res, err := TranslateBlockNumber("genesis")
	if err != nil {
		t.Errorf("TranslateBlockNumber #1: expected no error but got %v", err)
	}

	expected := int64(0)
	if res != expected {
		t.Errorf("TranslateBlockNumber #1: expected %v bug got %v", expected, res)
	}

	res, err = TranslateBlockNumber("latest")
	if err != nil {
		t.Errorf("TranslateBlockNumber #2: expected no error but got %v", err)
	}

	expected = int64(-1)
	if res != expected {
		t.Errorf("TranslateBlockNumber #2: expected %v bug got %v", expected, res)
	}

	res, err = TranslateBlockNumber("oldest")
	if err != nil {
		t.Errorf("TranslateBlockNumber #3: expected no error but got %v", err)
	}

	expected = int64(-2)
	if res != expected {
		t.Errorf("TranslateBlockNumber #3: expected %v bug got %v", expected, res)
	}

	res, err = TranslateBlockNumber("23679034")
	if err != nil {
		t.Errorf("TranslateBlockNumber #4: expected no error but got %v", err)
	}

	expected = int64(23679034)
	if res != expected {
		t.Errorf("TranslateBlockNumber #4: expected %v bug got %v", expected, res)
	}

	res, err = TranslateBlockNumber("unknown")
	if err == nil {
		t.Errorf("TranslateBlockNumber #5: expected an error")
	}

	expected = int64(0)
	if res != expected {
		t.Errorf("TranslateBlockNumber #5: expected %v bug got %v", expected, res)
	}
}

func TestParseStartingPosition(t *testing.T) {
	blockNumber, txIndex, err := ParseStartingPosition("genesis-0")
	if err != nil {
		t.Errorf("ParseStartingPosition #1: expected no error but got %v", err)
	}

	expected := [2]int64{0, 0}
	if blockNumber != expected[0] || txIndex != expected[1] {
		t.Errorf("ParseStartingPosition #1: expected %v bug got %v", expected, []int64{blockNumber, txIndex})
	}

	blockNumber, txIndex, err = ParseStartingPosition("latest-234454543")
	if err != nil {
		t.Errorf("ParseStartingPosition #2: expected no error but got %v", err)
	}

	expected = [2]int64{-1, 234454543}
	if blockNumber != expected[0] || txIndex != expected[1] {
		t.Errorf("ParseStartingPosition #2: expected %v bug got %v", expected, []int64{blockNumber, txIndex})
	}

	blockNumber, txIndex, err = ParseStartingPosition("oldest-12")
	if err != nil {
		t.Errorf("ParseStartingPosition #3: expected no error but got %v", err)
	}

	expected = [2]int64{-2, 12}
	if blockNumber != expected[0] || txIndex != expected[1] {
		t.Errorf("ParseStartingPosition #3: expected %v bug got %v", expected, []int64{blockNumber, txIndex})
	}

	blockNumber, txIndex, err = ParseStartingPosition("oldest")
	if err == nil {
		t.Errorf("ParseStartingPosition #4: expected an error")
	}

	expected = [2]int64{0, 0}
	if blockNumber != expected[0] || txIndex != expected[1] {
		t.Errorf("ParseStartingPosition #4: expected %v bug got %v", expected, []int64{blockNumber, txIndex})
	}

	blockNumber, txIndex, err = ParseStartingPosition("oldest-13a")
	if err == nil {
		t.Errorf("ParseStartingPosition #4: expected an error")
	}

	expected = [2]int64{0, 0}
	if blockNumber != expected[0] || txIndex != expected[1] {
		t.Errorf("ParseStartingPosition #4: expected %v bug got %v", expected, []int64{blockNumber, txIndex})
	}
}
