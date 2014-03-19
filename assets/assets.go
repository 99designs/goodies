package assets

import (
	"net/http"
	"strings"
)

type AssetPipeline interface {
	http.Handler
	AssetUrl(name string) (string, error)
	AssetContents(name string) ([]byte, error)
}

type prefixPipeline struct {
	AssetPipeline
	prefix string
}

// add the prefix when providing urls
func (s *prefixPipeline) AssetUrl(name string) (url string, err error) {
	url, err = s.AssetPipeline.AssetUrl(name)
	if err == nil {
		url = s.prefix + url
	}

	return
}

func (s *prefixPipeline) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/" + strings.TrimPrefix(r.URL.Path, s.prefix)
	s.AssetPipeline.ServeHTTP(w, r)
}

func Prefix(prefix string, p AssetPipeline) AssetPipeline {
	if len(prefix) == 0 {
		return p
	}
	return &prefixPipeline{
		prefix:        prefix,
		AssetPipeline: p,
	}
}
