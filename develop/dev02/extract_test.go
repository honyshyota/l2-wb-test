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
		isValid bool
	}{
		{
			name:    "with digits",
			val:     "a4bc2d5e",
			isValid: true,
		},
		{
			name:    "letters only",
			val:     "abcd",
			isValid: true,
		},
		{
			name:    "digits only",
			val:     "45",
			isValid: false,
		},
		{
			name:    "empty",
			val:     "",
			isValid: true,
		},
		{
			name:    "with shield 1",
			val:     `qwe\4\5`,
			isValid: true,
		},
		{
			name:    "with sield 2",
			val:     `qwe\45`,
			isValid: true,
		},
		{
			name:    "with shield 3",
			val:     `qwe\\5`,
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				_, err := extract.Extract(tc.val)
				assert.NoError(t, err)
			} else {
				_, err := extract.Extract(tc.val)
				assert.Error(t, err)
			}
		})
	}
}
