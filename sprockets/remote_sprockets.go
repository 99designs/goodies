package sprockets

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

type SprocketsServer struct {
	baseUrl      string
	manifestPath string
	assets       map[string]string
}

func NewSprocketsServer(baseUrl, manifestPath string) (ViewHelper, error) {
	s := SprocketsServer{baseUrl: baseUrl, manifestPath: manifestPath}
	err := s.loadManifest()
	return ViewHelper{&s}, err
}

func (s *SprocketsServer) GetAssetUrl(name string) (template.HTMLAttr, error) {
	src, ok := s.assets[name]
	if !ok {
		return template.HTMLAttr(s.baseUrl + "__NOT_FOUND__" + name), errors.New("Asset not found: " + name)
	}
	return template.HTMLAttr(s.baseUrl + src), nil
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

	return s.parseAssets(resp.Body)
}

func (s *SprocketsServer) parseAssets(r io.Reader) (map[string]string, error) {
	manifest := make(map[string]string)
	err := json.NewDecoder(r).Decode(&manifest)
	return manifest, err
}
