// Package sprockets integrates the power of sprockets ( http://getsprockets.org ) with your go program.
package sprockets

import (
	"io/ioutil"
	"net/http"

	"github.com/99designs/goodies/assets"
)

type SprocketsServer struct {
	internalUrl string
	externalUrl string
}

// NewSprocketsServer creates a Sprockets asset pipeline
func NewSprocketsServer(baseUrl string) (assets.AssetPipeline, error) {
	s := SprocketsServer{
		internalUrl: baseUrl,
		externalUrl: baseUrl,
	}

	return &s, nil
}

// NewPrivateSprocketsServer creates a Sprockets asset pipeline
// which uses a different baseUrl publicly
func NewPrivateSprocketsServer(externalUrl, internalUrl string) (assets.AssetPipeline, error) {
	s := SprocketsServer{
		internalUrl: internalUrl,
		externalUrl: externalUrl,
	}

	return &s, nil
}

func (s *SprocketsServer) GetAssetUrl(name string) (string, error) {
	return s.externalUrl + name, nil
}

func (s *SprocketsServer) GetAssetContents(name string) ([]byte, error) {
	resp, err := http.Get(s.internalUrl + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
