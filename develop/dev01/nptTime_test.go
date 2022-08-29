package ntptime_test

import (
	"fmt"
	"testing"

	ntpTime "github.com/honyshyota/l2-wb-test/develop/dev01"
	"github.com/stretchr/testify/assert"
)

func TestNptTime_NptTime(t *testing.T) {
	for i := 0; i < 5; i++ {
		time, err := ntpTime.NtpTime()
		assert.NoError(t, err)
		assert.NotEmpty(t, time)
		fmt.Println(time)
	}
}
