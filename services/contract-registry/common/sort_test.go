package common

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortStrings(t *testing.T) {
	tests := []struct {
		name string
		args []string
		res  []string
	}{
		{"base", []string{"z", "Z", "a", "A"}, []string{"A", "a", "Z", "z"}},
		{"opposite", []string{"Z", "z", "A", "a"}, []string{"A", "a", "Z", "z"}},
		{"bien", []string{"encore du travail", "1", "2", ".", "🛠"}, []string{".", "1", "2", "encore du travail", "🛠"}},
		{"empty", []string{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(Alphabetic(tt.args))
			assert.Equal(t, tt.res, tt.args)
		})
	}
}
