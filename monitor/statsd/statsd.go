// Package statsd records the time taken to respond to HTTP requests in statsd
package statsd

import (
	"github.com/peterbourgon/g2s"
	"log"
)

func Dial(proto, addr string) *Statsd {
	st, err := g2s.Dial(proto, addr)

	if err != nil {
		log.Printf("Couldn't initiate statsd with address '%s' because %s", addr, err)
		return &Statsd{g2s.Noop()}
	}
	return &Statsd{st}
}

type Statsd struct {
	g2s.Statter
}
