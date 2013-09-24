package goanna

import (
	. "launchpad.net/gocheck"
	"net/http"
	"testing"
	"time"
)

type HandlerSuite struct {
	TestController
}

type TestController struct {
	*Controller
	Response
}

var _ = Suite(&HandlerSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *HandlerSuite) TestCreateLoginToken(c *C) {
	req := Must(http.NewRequest("GET", "/foo", nil))

	factory := func() ControllerInterface { return &s.TestController }
	h := NewHandler(factory, "ControllerMethod")

	s.TestController.Controller = NewController(req, nil)
	s.TestController.Response = NewResponse()
	resp := h.getResponse(req)
	c.Check(resp, Equals, s.TestController.Response)
}

func (s *TestController) Init(*http.Request)         {}
func (s *TestController) Session() Session           { return NopSession{} }
func (s *TestController) ControllerMethod() Response { return s.Response }

func Must(r *http.Request, e error) *http.Request {
	if e != nil {
		panic(e)
	}
	return r
}

type NopSession struct {
}

func (s NopSession) GetId() string              { return "" }
func (s NopSession) Get(_ string) string        { return "" }
func (s NopSession) Set(_string, _ string)      {}
func (s NopSession) SetMaxAge(_ time.Duration)  {}
func (s NopSession) HasExpired() bool           { return false }
func (s NopSession) Clear()                     {}
func (s NopSession) WriteToResponse(_ Response) {}
