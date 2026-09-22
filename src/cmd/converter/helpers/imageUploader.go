package helpers

import (
	"bytes"
	"path/filepath"

	"github.com/3l4l5/blog/src/cmd/converter/core/images"
	"github.com/3l4l5/blog/src/lib"
)

type ImageUploader func(imageArray []images.Image) ([]images.ImageUploaded, error)

func ImageUploaderFactory(option string) ImageUploader {
	switch option {
	case "publish":
		return s3ImageUploader
	case "build":
		return localImageUploader
	default:
		return nil
	}
}

func s3ImageUploader(inputImages []images.Image) ([]images.ImageUploaded, error) {
	client, err := lib.NewS3Manager()
	if err != nil {
		return nil, err
	}

	env, err := lib.GetEnv()
	if err != nil {
		return nil, err
	}

	result := []images.ImageUploaded{}

	for _, img := range inputImages {
		imgByte, err := lib.GetFilesAsByte(img.OriginPath)
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(imgByte)
		path := filepath.Join(img.ArticleId, filepath.Base(img.OriginPath))
		url, err := client.Upload(env.S3_BLOG_IMAGE_BUCKET_NAME, path, reader)
		if err != nil {
			return nil, err
		}
		result = append(result, images.ImageUploaded{
			ArticleId:   img.ArticleId,
			OriginPath:  img.OriginPath,
			UploadedUrl: url,
		})
	}
	return result, nil
}

func localImageUploader(inputImages []images.Image) ([]images.ImageUploaded, error) {
	result := []images.ImageUploaded{}
	for _, img := range inputImages {
		imgByte, err := lib.GetFilesAsByte(img.OriginPath)
		if err != nil {
			return []images.ImageUploaded{}, nil
		}
		reader := bytes.NewReader(imgByte)
		err = lib.Save("dist/images/"+img.OriginPath, reader)
		if err != nil {
			return nil, err
		}
		result = append(result,
			images.ImageUploaded{
				ArticleId:   img.ArticleId,
				OriginPath:  img.OriginPath,
				UploadedUrl: "/images/" + img.OriginPath,
			})
	}
	return result, nil
}
