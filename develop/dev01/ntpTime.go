package ntpTime

import (
	"time"

	"github.com/beevik/ntp"
)

func NtpTime() (time.Time, error) {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return time, err
	}

	return time, err
}
