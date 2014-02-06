package assets

type NoopAssetPipeline struct {
}

func (s NoopAssetPipeline) GetAssetUrl(name string) (a string, b error) {
	return
}

func (s NoopAssetPipeline) GetAssetContents(name string) (a []byte, b error) {
	return
}

func NoopAssetHelper() AssetHelper {
	return AssetHelper{NoopAssetPipeline{}}
}
