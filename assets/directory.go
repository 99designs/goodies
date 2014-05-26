package assets

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// AssetDirectory is an AssetPipeline for a
// directory containing asssets
type AssetDirectory struct {
	fileserver http.Handler
	path       string
}

func NewAssetDirectory(path string) (AssetPipeline, error) {
	s := AssetDirectory{
		fileserver: http.FileServer(http.Dir(path)),
		path:       path,
	}

	return &s, nil
}

func (s *AssetDirectory) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" || r.URL.Path == "/" {
		http.NotFound(w, r)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/")

	if filename[0] != '/' {
		filename = "/" + filename
	}
	r.URL.Path = filename

	s.fileserver.ServeHTTP(w, r)
}

func (s *AssetDirectory) AssetUrl(name string) (string, error) {
	return name, nil
}

func (s *AssetDirectory) AssetContents(filename string) (content []byte, err error) {
	return ioutil.ReadFile(s.path + "/" + filename)
}
