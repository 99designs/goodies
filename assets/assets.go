package assets

import "net/http"

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

// Prefix returns an AssetPipeline that prefixes URLs with the
// given string
func Prefix(prefix string, p AssetPipeline) AssetPipeline {
	if prefix == "" {
		return p
	}

	return &prefixPipeline{
		prefix:        prefix,
		AssetPipeline: p,
	}
}
