package assets

import (
	"io/ioutil"
	"net/http"
	"path"
)

// AssetDirectory is an AssetPipeline for a
// directory containing asssets
type AssetDirectory struct {
	fileserver http.Handler
	path       string
}

func NewAssetDirectory(path string) *AssetDirectory {
	return &AssetDirectory{
		fileserver: http.FileServer(http.Dir(path)),
		path:       path,
	}
}

func (s *AssetDirectory) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" || r.URL.Path == "/" {
		http.NotFound(w, r) // avoid showing the index page that FileServer throws up
	} else {
		s.fileserver.ServeHTTP(w, r)
	}
}

func (s *AssetDirectory) AssetUrl(name string) (string, error) {
	return name, nil
}

func (s *AssetDirectory) AssetContents(filename string) (content []byte, err error) {
	return ioutil.ReadFile(path.Join(s.path, filename))
}
