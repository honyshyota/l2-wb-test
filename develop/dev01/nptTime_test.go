package ntpTime

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNptTime_NptTime(t *testing.T) {
	for i := 0; i < 10; i++ {
		time, err := NtpTime()
		assert.NoError(t, err)
		assert.NotEmpty(t, time)
		fmt.Println(time)
	}
}
