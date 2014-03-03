package sprockets

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/goodies/assets"
)

// CompiledAssetDirectory is an AssetPipeline for a
// directory containing compiled assets and a sprockets
// manifest.json file
type CompiledAssetDirectory struct {
	fileserver http.Handler
	path       string
	manifest   map[string]string
}

func NewCompiledAssetDirectory(path string) (assets.AssetPipeline, error) {
	s := CompiledAssetDirectory{
		fileserver: http.FileServer(http.Dir(path)),
		path:       path,
	}
	err := s.loadManifest()

	return &s, err
}

func (s *CompiledAssetDirectory) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" || r.URL.Path == "/" {
		http.NotFound(w, r)
		return
	}

	filename, ok := s.manifest[strings.TrimPrefix(r.URL.Path, "/")]
	if ok {
		if filename[0] != '/' {
			filename = "/" + filename
		}
		r.URL.Path = filename
	}
	s.fileserver.ServeHTTP(w, r)
}

func (s *CompiledAssetDirectory) AssetUrl(name string) (url string, err error) {
	url, ok := s.manifest[name]
	if !ok {
		err = errors.New("Asset not found: " + name)
	}

	return
}

func (s *CompiledAssetDirectory) AssetContents(name string) (content []byte, err error) {
	filename, ok := s.manifest[name]
	if !ok {
		err = errors.New("Asset not found: " + name)
	} else {
		content, err = ioutil.ReadFile(s.path + "/" + filename)
	}

	return
}

func (s *CompiledAssetDirectory) loadManifest() error {
	file, err := os.Open(s.path + "/manifest.json")
	if err != nil {
		return err
	}
	manifest := Manifest{
		Assets: make(map[string]string),
	}
	err = json.NewDecoder(file).Decode(&manifest)

	s.manifest = manifest.Assets

	return err
}

type Manifest struct {
	Assets map[string]string `json:"assets"`
}
