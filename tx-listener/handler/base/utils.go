package base

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	positionRegexp  = `^(?P<chain>[0-9]+):(?P<blockNumber>genesis|latest|oldest|\d+)(-(?P<txIndex>\d+))?`
	positionPattern = regexp.MustCompile(positionRegexp)
)

// ParseStartingPosition extract chainID, blockNumber and TxIndex from a formatted starting position string
func ParsePosition(position string) (string, *Position, error) {
	match := positionPattern.FindStringSubmatch(position)
	if len(match) != 5 {
		return "", nil, fmt.Errorf("could not parse position %q (expected format: %q)", position, positionRegexp)
	}

	blockNumber, err := ParseBlock(match[2])
	if err != nil {
		return "", nil, fmt.Errorf("could not parse position %q: %v", position, err)
	}

	if match[4] == "" {
		return match[1], &Position{blockNumber, 0}, nil
	}

	txIndex, err := strconv.ParseInt(match[4], 10, 64)
	if err != nil {
		return "", nil, fmt.Errorf("could not parse position %q: %v", position, err)
	}

	return match[1], &Position{blockNumber, txIndex}, nil
}

var (
	genesis = "genesis"
	latest  = "latest"
	oldest  = "oldest"
)

// ParseBlock translate a starting block number into its integer value
func ParseBlock(block string) (int64, error) {
	switch block {
	case genesis:
		return 0, nil
	case latest:
		return -1, nil
	case oldest:
		return -2, nil
	default:
		res, err := strconv.ParseInt(block, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%q is an invalid block (should be one of %q or an integer)", block, []string{genesis, oldest, latest})
		}
		return res, nil
	}
}
