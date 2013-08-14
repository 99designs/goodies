package http

// Type HttpError is used with PanicHandler as
// a convenient way to be able to panic and
// ensure that a particular http status code will
// be returned.
type HttpError struct {
	Error      error // the original error
	StatusCode int   // a suggested http error code
}
