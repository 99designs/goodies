package http

// HttpError is used as a convenient way to be able to panic and
// ensure that a particular http status code will be returned
type HttpError interface {
	Error() string   // the original error
	StatusCode() int // a suggested http error code
}

type httpError struct {
	err error // the original error
	sc  int   // a suggested http error code
}

func (this httpError) Error() string {
	return this.err.Error()
}

func (this httpError) StatusCode() int {
	return this.sc
}

func NewHttpError(err error, sc int) HttpError {
	return &httpError{err, sc}
}
