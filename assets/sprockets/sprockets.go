// Package sprockets integrates the power of sprockets ( http://getsprockets.org ) with your go program.
package sprockets

import (
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// SprocketsServer is an implementation of AssetPipeline that
// uses an external Sprockets server for compiling assets
// on demand
type SprocketsServer struct {
	http.Handler
	target *url.URL
}

// NewSprocketsServer creates a Sprockets asset pipeline
func NewSprocketsServer(target *url.URL) (*SprocketsServer, error) {
	s := SprocketsServer{
		Handler: httputil.NewSingleHostReverseProxy(target),
		target:  target,
	}

	return &s, nil
}

func (s *SprocketsServer) AssetUrl(name string) (string, error) {
	return name, nil
}

func (s *SprocketsServer) AssetContents(name string) ([]byte, error) {
	url := s.target.String() + name
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
