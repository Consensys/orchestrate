package memory

import (
	"fmt"
)

func computeKey(key, suffix string) string {
	return fmt.Sprintf("%v-%v", key, suffix)
}
