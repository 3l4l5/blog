package images

type Image struct {
	ArticleId  string
	OriginPath string
}

type ImageUploaded struct {
	ArticleId   string
	OriginPath  string
	UploadedUrl string
}
