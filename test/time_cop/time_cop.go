// Handy testing stub on time.Time
package time_cop

import (
	"time"
)

func Freeze(t int) func() time.Time {
	return func() time.Time {
		return time.Unix(int64(t), 0)
	}
}
