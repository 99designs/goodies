package prereader

import (
	"errors"
	"io"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: Implement a fake IO reader that returns an error when you try to read from it, then see what happens.

type PrereaderSuite struct {
	server *httptest.Server
}

var _ = Suite(&PrereaderSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *PrereaderSuite) SetUpTest(c *C) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Testing 1 2 3"))
	})

	s.server = httptest.NewServer(Decorate(h))
}

func (s *PrereaderSuite) TearDownTest(c *C) {
	s.server.Close()
}

func (s *PrereaderSuite) request(body io.Reader) *http.Request {
	req, err := http.NewRequest("GET", s.server.URL+"/foo/bar", body)
	panicIf(err)
	return req
}

func (s *PrereaderSuite) TestPrereadBody(c *C) {
	req := s.request(TimeoutReader{})
	err := PrereadBody(req)
	c.Check(err, Equals, TimeoutError)
}

func (s *PrereaderSuite) TestDecorator(c *C) {
	req := s.request(TimeoutReader{})
	rw := httptest.NewRecorder()
	PrereadHandler{}.ServeHttp(rw, req)
	c.Check(rw.Code, Equals, http.StatusRequestTimeout)
	c.Check(string(rw.Body.Bytes()), Equals, "Error reading request: fake i/o timeout\n")
}

type TimeoutReader struct{}

var TimeoutError = errors.New("fake i/o timeout")

func (s TimeoutReader) Read(p []byte) (n int, err error) {
	panic(TimeoutError)
}

func (s TimeoutReader) Close() (err error) {
	return nil
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
