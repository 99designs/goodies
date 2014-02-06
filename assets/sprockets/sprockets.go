// Package sprockets integrates the power of sprockets ( http://getsprockets.org ) with your go program.
package sprockets

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/99designs/goodies/assets"
)

type SprocketsServer struct {
	baseUrl      string
	manifestPath string
	assets       map[string]string
}

func NewSprocketsServer(baseUrl, manifestPath string) (assets.AssetPipeline, error) {
	s := SprocketsServer{baseUrl: baseUrl, manifestPath: manifestPath}
	err := s.loadManifest()

	return &s, err
}

func (s *SprocketsServer) GetAssetUrl(name string) (string, error) {
	src, ok := s.assets[name]
	if !ok {
		return (s.baseUrl + "__NOT_FOUND__" + name), errors.New("Asset not found: " + name)
	}
	return (s.baseUrl + src), nil
}

func (s *SprocketsServer) GetAssetContents(name string) ([]byte, error) {
	resp, err := http.Get(s.baseUrl + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (s *SprocketsServer) loadManifest() error {
	var err error
	if s.assets == nil {
		s.assets, err = s.fetchAssetList(s.baseUrl + s.manifestPath)
	}
	return err
}

func (s *SprocketsServer) fetchAssetList(url string) (map[string]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Couldn't read assets from " + url)
	}

	manifest := make(map[string]string)
	err = json.NewDecoder(resp.Body).Decode(&manifest)

	return manifest, err
}
