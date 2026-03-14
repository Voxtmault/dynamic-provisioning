package object

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const presignedExpiry = 30 * time.Minute

type S3Client struct {
	client *minio.Client
	bucket string
}

func NewS3Client(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*S3Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 client: %w", err)
	}

	return &S3Client{client: client, bucket: bucket}, nil
}

func (s *S3Client) GeneratePresignedURL(objectKey string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)

	presignedURL, err := s.client.PresignedGetObject(
		context.Background(), s.bucket, objectKey, expiry, reqParams,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}

	return presignedURL.String(), nil
}

func (s *S3Client) GeneratePresignedPutURL(objectKey string, expiry time.Duration) (string, error) {
	presignedURL, err := s.client.PresignedPutObject(
		context.Background(), s.bucket, objectKey, expiry,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned put url: %w", err)
	}

	return presignedURL.String(), nil
}
