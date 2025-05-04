package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

type Client struct {
	s3Client   *s3.Client
	s3Uploader *manager.Uploader
	bucketName string
}

func NewClient(accessKeyId, accessKeySecret, endpoint string, bucket string) (*Client, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: endpoint,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)

	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(s3Client)

	return &Client{
		s3Client:   s3Client,
		s3Uploader: uploader,
		bucketName: bucket,
	}, nil
}

func (c *Client) UploadFile(ctx context.Context, key string, body io.Reader, contentType string) error {
	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	}

	_, err := c.s3Uploader.Upload(ctx, uploadInput)
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	return nil
}
