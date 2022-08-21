package anagramSearch_test

import (
	"testing"

	anagramSearch "github.com/honyshyota/l2-wb-test/develop/dev04"
	"github.com/stretchr/testify/assert"
)

func TestExtract_Extract(t *testing.T) {
	testCases := []struct {
		name  string
		val   []string
		equal map[string][]string
	}{
		{
			name: "valid",
			val:  []string{"листок", "пятак", "пятка", "слиток", "столик", "тяпка", "кот", "ток"},
			equal: map[string][]string{
				"кот":    {"кот", "ток"},
				"листок": {"листок", "слиток", "столик"},
				"пятак":  {"пятак", "пятка", "тяпка"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val := anagramSearch.FindAnagrams(&tc.val)
			assert.Equal(t, val, &tc.equal)
		})
	}
}
