// Package timecop implements a stub to use in place of time.Time
package time_cop

import (
	"time"
)

func Freeze(t int) func() time.Time {
	return func() time.Time {
		return time.Unix(int64(t), 0)
	}
}
