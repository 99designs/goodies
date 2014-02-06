package sprockets

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/99designs/goodies/assets"
)

type CompiledAssetDirectory struct {
	dirpath string
	prefix  string
	assets  map[string]string
}

func NewCompiledAssetDirectory(dirpath, prefix string) (assets.AssetPipeline, error) {
	s := CompiledAssetDirectory{dirpath: dirpath, prefix: prefix}
	err := s.loadManifest()

	return &s, err
}

func (s *CompiledAssetDirectory) GetAssetUrl(name string) (string, error) {
	src, ok := s.assets[name]
	if !ok {
		return (s.prefix + "__NOT_FOUND__" + name), errors.New("Asset not found: " + name)
	}
	return (s.prefix + src), nil
}

func (s *CompiledAssetDirectory) GetAssetContents(name string) ([]byte, error) {
	return ioutil.ReadFile(s.dirpath + "/" + s.assets[name])
}

func (s *CompiledAssetDirectory) loadManifest() (err error) {
	var file *os.File
	file, err = os.Open(s.dirpath + "/manifest.json")
	if err != nil {
		return
	}
	manifest := Manifest{Assets: make(map[string]string)}
	err = json.NewDecoder(file).Decode(&manifest)

	s.assets = manifest.Assets

	return err
}

type Manifest struct {
	Assets map[string]string `json:"assets"`
}
