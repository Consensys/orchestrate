// +build unit

package common

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunInParallel(t *testing.T) {
	executionDone := make(chan int, 2)

	InParallel(
		func() {
			executionDone <- 1
		},
		func() {
			executionDone <- 2
		},
	)

	// Wait for both functions to be executed
	var res = []int{
		<-executionDone,
		<-executionDone,
	}

	// Results can come in any order
	sort.Ints(res)

	assert.Equal(t, []int{1, 2}, res, "Should collect result values from all functions")
}
