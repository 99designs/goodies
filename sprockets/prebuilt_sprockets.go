package sprockets

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"os"
)

type SprocketsDirectory struct {
	dirpath string
	prefix  string
	assets  map[string]string
}

func NewSprocketsDirectory(dirpath, prefix string) (ViewHelper, error) {
	s := &SprocketsDirectory{dirpath: dirpath, prefix: prefix}
	err := s.loadManifest()
	return ViewHelper{s}, err
}

func (s *SprocketsDirectory) GetAssetUrl(name string) (template.HTMLAttr, error) {
	src, ok := s.assets[name]
	if !ok {
		return template.HTMLAttr(s.prefix + "__NOT_FOUND__" + name), errors.New("Asset not found: " + name)
	}
	return template.HTMLAttr(s.prefix + src), nil
}

func (s *SprocketsDirectory) GetAssetContents(name string) ([]byte, error) {
	return ioutil.ReadFile(s.dirpath + "/" + s.assets[name])
}

func (s *SprocketsDirectory) loadManifest() (err error) {
	var file *os.File
	if s.assets == nil {
		file, err = os.Open(s.dirpath + "/manifest.json")
		if err != nil {
			return
		}
		s.assets, err = s.parseManifest(file)
	}
	return err
}

func (s *SprocketsDirectory) parseManifest(r io.Reader) (map[string]string, error) {
	manifest := Manifest{Assets: make(map[string]string)}
	err := json.NewDecoder(r).Decode(&manifest)
	return manifest.Assets, err
}

type Manifest struct {
	Assets map[string]string `json:"assets"`
}
