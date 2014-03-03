package assets

import "net/http"

type noopAssetPipeline struct {
	http.Handler
}

func (s noopAssetPipeline) AssetUrl(name string) (a string, b error) {
	return
}

func (s noopAssetPipeline) AssetContents(name string) (a []byte, b error) {
	return
}

var NoopAssetPipeline = &noopAssetPipeline{
	Handler: http.NotFoundHandler(),
}
