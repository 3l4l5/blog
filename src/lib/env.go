package lib

import (
	"fmt"
	"os"
)

type Env struct {
	AWS_ACCESS_KEY            string
	AWS_ACCESS_SECRET         string
	AWS_REGION                string
	AWS_ENDPOINT              string
	S3_BLOG_IMAGE_BUCKET_NAME string
}

func GetEnv() (*Env, error) {
	key := os.Getenv("AWS_ACCESS_KEY")
	secret := os.Getenv("AWS_ACCESS_SECRET")
	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("AWS_ENDPOINT")
	bucket := os.Getenv("S3_BLOG_IMAGE_BUCKET_NAME")
	if key == "" ||
		secret == "" ||
		region == "" ||
		endpoint == "" ||
		bucket == "" {
		return &Env{}, fmt.Errorf("Env value is invalid.")
	}
	return &Env{
		AWS_ACCESS_KEY:            key,
		AWS_ACCESS_SECRET:         secret,
		AWS_REGION:                region,
		AWS_ENDPOINT:              endpoint,
		S3_BLOG_IMAGE_BUCKET_NAME: bucket,
	}, nil
}
