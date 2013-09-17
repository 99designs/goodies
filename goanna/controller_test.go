package goanna

import (
	"bytes"
	"fmt"
	. "launchpad.net/gocheck"
	"log"
	"net/http"
	"strings"
	"testing"
)

type ControllerSuite struct{}

var _ = Suite(&ControllerSuite{})

func TestErrorLogging(t *testing.T) { TestingT(t) }

func (s *ControllerSuite) TestErrorLogging(c *C) {
	req, _ := http.NewRequest("GET", "/", nil)
	con := NewController(req, nil)
	output := bytes.Buffer{}
	con.logger = log.New(&output, "", 0)
	con.LogRequest("Just for testing")
	fmt.Print(string(output.Bytes()))

	startOfLog := `
*** Diagnostic Log ***
Reason for diagnostic: Just for testing
Url: /
Method: GET
Timestamp:`
	out := string(output.Bytes())
	c.Check(out, HasPrefix, startOfLog)
	// Ensure stack trace is printed
	c.Check(out, Contains, "goanna/controller.go")
	c.Check(out, Contains, "goanna/controller_test.go")
}

type SimpleChecker struct {
	*CheckerInfo
	CheckFn func([]interface{}, []string) (bool, string)
}

func (h SimpleChecker) Info() *CheckerInfo {
	return h.CheckerInfo
}

func (h SimpleChecker) Check(params []interface{}, names []string) (result bool, error string) {
	return h.CheckFn(params, names)
}

var HasPrefix = SimpleChecker{&CheckerInfo{Name: "HasPrefix", Params: []string{"obtained", "expected"}}, CheckHasPrefix}

func CheckHasPrefix(params []interface{}, names []string) (result bool, error string) {
	s := params[0].(string)
	pre := params[1].(string)
	return strings.HasPrefix(s, pre), ""
}

var Contains = SimpleChecker{&CheckerInfo{Name: "HasPrefix", Params: []string{"obtained", "expected"}}, CheckContains}

func CheckContains(params []interface{}, names []string) (result bool, error string) {
	s := params[0].(string)
	pre := params[1].(string)
	return strings.Contains(s, pre), ""
}
