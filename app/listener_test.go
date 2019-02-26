package app

import "testing"

func TestParseStartingPosition(t *testing.T) {
	chain, position, err := ParseStartingPosition("0xab:genesis")
	if chain != "0xab" || position.BlockNumber != 0 || position.TxIndex != 0 || err != nil {
		t.Errorf("ParseStartingPosition #1: got %q %v %v %v", chain, position.BlockNumber, position.TxIndex, err)
	}

	chain, position, err = ParseStartingPosition("0xb:24-124")
	if chain != "0xb" || position.BlockNumber != 24 || position.TxIndex != 124 || err != nil {
		t.Errorf("ParseStartingPosition #2: got %q %v %v %v", chain, position.BlockNumber, position.TxIndex, err)
	}

	chain, position, err = ParseStartingPosition("0x3:latest-0")
	if chain != "0x3" || position.BlockNumber != -1 || position.TxIndex != 0 || err != nil {
		t.Errorf("ParseStartingPosition #3: got %q %v %v %v", chain, position.BlockNumber, position.TxIndex, err)
	}

	chain, position, err = ParseStartingPosition("3:latest-0")
	if chain != "" || position != nil || err == nil {
		t.Errorf("ParseStartingPosition #4: got %q %v %v %v", chain, position.BlockNumber, position.TxIndex, err)
	}
}

func TestParseStartingPositions(t *testing.T) {
	positions := []string{
		"0xab:genesis",
		"0xb:24-124",
		"0x3:latest-0",
	}

	parsedPositions, err := ParseStartingPositions(positions)
	if err != nil {
		t.Errorf("ParseStartingPositions #1: expected no error but got %v", err)
	} else {
		if len(parsedPositions) != 3 {
			t.Errorf("ParseStartingPositions #2: expected %v positions but got %v", 3, len(parsedPositions))
		} else {
			if parsedPositions["0xab"].BlockNumber != 0 || parsedPositions["0xab"].TxIndex != 0 {
				t.Errorf("ParseStartingPosition #3: unexpected position")
			}

			if parsedPositions["0xb"].BlockNumber != 24 || parsedPositions["0xb"].TxIndex != 124 {
				t.Errorf("ParseStartingPosition #4: unexpected position")
			}
		}
	}
}
