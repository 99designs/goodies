// Package goodies/expvar exposes several commonly-used variables
// via the built-in expvar
package expvar

import (
	"expvar" // Adds '/debug/vars' http route
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type exposedString struct {
	str string
}

func (e exposedString) String() string {
	return fmt.Sprintf("%#v", e.str)
}

func init() {
	expvar.NewInt("NumCPUs").Set(int64(runtime.NumCPU()))

	revision, err := exec.Command("git", "log", "-1", "--pretty=oneline", "HEAD").Output()
	if err != nil {
		expvar.NewString("revision").Set(fmt.Sprintf("Could not determine git version: %s", err))
	} else {
		expvar.NewString("revision").Set(strings.TrimSpace(string(revision)))
	}

	env := expvar.NewMap("env")
	for _, val := range os.Environ() {
		parts := strings.SplitN(val, "=", 2)
		if len(parts) >= 2 {
			env.Set(parts[0], exposedString{parts[1]})
		}
	}
}
