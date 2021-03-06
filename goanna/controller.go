// Package goanna is an MVC toolkit with routing, controllers and sessions
package goanna

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
)

// ControllerInterface is implemented by controllers
type ControllerInterface interface {
	Init()
	Session() Session
	SetRequest(*http.Request)
}

// Controller is an embeddable type for controllers
type Controller struct {
	Request       *Request
	sessionFinder SessionFinder
}

// SetRequest injects a request into the controller
func (c *Controller) SetRequest(req *http.Request) {
	c.Request = &Request{Request: req}
}

// Session returns the session for the current request
func (c *Controller) Session() Session {
	if c.Request.session == nil {
		c.Request.session = c.sessionFinder(c.Request)
	}
	return c.Request.session
}

// NewController instantiates a new Controller using the given request and session store
func NewController(sessionFinder SessionFinder) Controller {
	return Controller{
		sessionFinder: sessionFinder,
	}
}

// LogRequest dumps the current request to stdout
func (c *Controller) LogRequest(reason string) {
	LogRequest(c.Request, reason)
}

// RenderView renders a template string using the provided template and vars struct
// and returns the rendered tamplate
func (c *Controller) RenderView(templateStr string, vars interface{}) []byte {
	t, err := template.New("RenderView").Parse(templateStr)
	if err != nil {
		panic(err.Error())
	}
	return c.RenderTemplate(t, vars)
}

// RenderView renders a template using the provided template and vars struct
// and returns the rendered tamplate
func RenderTemplate(t *template.Template, vars interface{}) []byte {
	out := bytes.NewBuffer(nil)
	err := t.Execute(out, vars)
	if err != nil {
		panic(err.Error())
	}
	return out.Bytes()
}
func (c *Controller) RenderTemplate(t *template.Template, vars interface{}) []byte {
	return RenderTemplate(t, vars)
}

// Render renders a template using the provided template and vars struct
// and returns a response with the rendered template
func RenderTemplateResponse(t *template.Template, vars interface{}) *OkResponse {
	response := NewResponse()
	response.SetContent(RenderTemplate(t, vars))
	return response
}
func (c *Controller) RenderTemplateResponse(t *template.Template, vars interface{}) *OkResponse {
	return RenderTemplateResponse(t, vars)
}

// Redirect returns a new RedirectResponse
// (HTTP 302)
func Redirect(urlStr string) *RedirectResponse {
	return NewRedirectResponse(urlStr)
}

func (c *Controller) Redirect(urlStr string) *RedirectResponse {
	return NewRedirectResponse(urlStr)
}

// PermanentRedirect returns a new RedirectResponse with a permanent redirect
// (HTTP 301)
func (c *Controller) PermanentRedirect(urlStr string) *RedirectResponse {
	return NewPermanentRedirectResponse(urlStr)
}

// RedirectRoute returns a RedirectResponse to the route
func (c *Controller) RedirectRoute(routeName string, args ...string) *RedirectResponse {
	return NewRedirectResponse(c.UrlFor(routeName, args...).String())
}

// Render renders a template using the provided template and vars struct
// and returns a response with the rendered template
func (c *Controller) Render(templateStr string, vars interface{}) *OkResponse {
	response := NewResponse()
	response.SetContent(c.RenderView(templateStr, vars))
	return response
}

// UrlFor is helper function for controllers
func (c *Controller) UrlFor(routeName string, args ...string) *url.URL {
	return UrlFor(routeName, args...)
}
