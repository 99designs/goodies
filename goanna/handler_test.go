package goanna

import (
	"net/http"
	"testing"
	. "gopkg.in/check.v1"
)

type HandlerSuite struct {
	TestController
}

type TestController struct {
	Controller
	Response
}

var _ = Suite(&HandlerSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *HandlerSuite) TestCreateLoginToken(c *C) {
	req := Must(http.NewRequest("GET", "/foo", nil))

	factory := func() ControllerInterface { return &s.TestController }
	h := NewHandler(factory, "ControllerMethod")

	s.TestController.Controller = NewController(nil)
	s.TestController.Controller.SetRequest(req)
	s.TestController.Response = NewResponse()
	resp := h.getResponse(req)
	c.Check(resp, Equals, s.TestController.Response)
}

func (s *TestController) Init()                      {}
func (s *TestController) SetRequest(*http.Request)   {}
func (s *TestController) Session() Session           { return NopSession{} }
func (s *TestController) ControllerMethod() Response { return s.Response }

func Must(r *http.Request, e error) *http.Request {
	if e != nil {
		panic(e)
	}
	return r
}
