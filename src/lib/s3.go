package lib

import (
	"context"
	"io"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Manager struct {
	client *s3.Client
}

func NewS3Manager() (*S3Manager, error) {
	env, err := GetEnv()
	if err != nil {
		return nil, err
	}
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(env.AWS_REGION),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				env.AWS_ACCESS_KEY,
				env.AWS_ACCESS_SECRET,
				"",
			),
		),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(env.AWS_ENDPOINT)
		o.UsePathStyle = true
	})

	return &S3Manager{
		client: client,
	}, nil
}

func (s *S3Manager) Upload(bucket string, path string, file io.Reader) (string, error) {
	env, err := GetEnv()
	if err != nil {
		return "", err
	}
	_, err = s.client.PutObject(
		context.Background(),
		&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(path),
			Body:   file,
		},
	)
	if err != nil {
		return "", err
	}
	joined, err := url.JoinPath(
		env.AWS_ENDPOINT,
		bucket,
		path,
	)
	if err != nil {
		return "", err
	}
	return joined, nil
}

func (s *S3Manager) Exist(bucket string, path string) (bool, error) {
	_, err := s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err == nil {
		return true, nil
	}

	return false, err
}
