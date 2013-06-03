// Null Value wrapper (for github.com/peterbourgon/g2s).
package statsd

import (
	"github.com/peterbourgon/g2s"
	"log"
	"time"
)

func Dial(proto, addr string) *Statsd {
	st, err := g2s.Dial(proto, addr)

	if err != nil {
		log.Printf("Couldn't initiate statsd with address: '%s'", addr)
		return nil
	}
	return &Statsd{st}
}

type Statsd struct {
	g2s.Statter
}

func (sd Statsd) Counter(sampleRate float32, bucket string, n ...int) {
	if sd.Statter != nil {
		sd.Statter.Counter(sampleRate, bucket, n...)
	}
}

func (sd Statsd) Timing(sampleRate float32, bucket string, d ...time.Duration) {
	if sd.Statter != nil {
		sd.Statter.Timing(sampleRate, bucket, d...)
	}
}

func (sd Statsd) Gauge(sampleRate float32, bucket string, value ...string) {
	if sd.Statter != nil {
		sd.Gauge(sampleRate, bucket, value...)
	}
}
