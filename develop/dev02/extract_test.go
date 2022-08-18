package extract_test

import (
	"testing"

	extract "github.com/honyshyota/l2-wb-test/develop/dev02"
	"github.com/stretchr/testify/assert"
)

func TestExtract_Extract(t *testing.T) {
	testCases := []struct {
		name    string
		val     string
		equal   string
		isValid bool
	}{
		{
			name:    "with digits",
			val:     "a4bc2d5e",
			equal:   "aaaabccddddde",
			isValid: true,
		},
		{
			name:    "letters only",
			val:     "abcd",
			equal:   "abcd",
			isValid: true,
		},
		{
			name:    "digits only",
			isValid: true,
		},
		{
			name:    "digits only",
			val:     "45",
			equal:   "",
			isValid: false,
		},
		{
			name:    "empty",
			val:     "",
			equal:   "",
			isValid: true,
		},
		{
			name:    "with shield 1",
			val:     `qwe\4\5`,
			equal:   "qwe45",
			isValid: true,
		},
		{
			name:    "with sield 2",
			val:     `qwe\45`,
			equal:   "qwe44444",
			isValid: true,
		},
		{
			name:    "with shield 3",
			val:     `qwe\\5`,
			equal:   `qwe\\\\\`,
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				val, err := extract.Extract(tc.val)
				assert.Equal(t, val, tc.equal)
				assert.NoError(t, err)
			} else {
				val, err := extract.Extract(tc.val)
				assert.Equal(t, tc.equal, val)
				assert.Error(t, err)
			}
		})
	}
}
