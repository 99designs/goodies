package assets

import (
	"net/http"
	"strings"
)

type PackagedAsset struct {
	f func(name string) ([]byte, error)
}

func NewPackagedAsset(f func(name string) ([]byte, error)) AssetPipeline {
	pa := &PackagedAsset{f}
	return pa
}

func (p *PackagedAsset) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" || r.URL.Path == "/" {
		http.NotFound(w, r)
		return
	}

	// filename similar to "public/my-file.html"

	filename := strings.TrimPrefix(r.URL.Path, "/")
	buffer, err := p.f(filename)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Write(buffer)
}

func (s *PackagedAsset) AssetUrl(name string) (string, error) {
	return name, nil
}

func (s *PackagedAsset) AssetContents(filename string) (content []byte, err error) {
	// filename parameter is of format "public/something.html" - just what we want
	return s.f(filename)
}
