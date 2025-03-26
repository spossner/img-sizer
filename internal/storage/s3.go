package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client handles S3 operations
type S3Client struct {
	client *s3.Client
}

// NewS3Client creates a new S3 client using environment credentials
func NewS3Client(ctx context.Context) (*S3Client, error) {
	// Get AWS configuration from environment variables
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Validate required environment variables
	if accessKeyID == "" || secretAccessKey == "" {
		return nil, fmt.Errorf("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY must be set")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "eu-central-1" // Default region if not set
	}

	creds := credentials.NewStaticCredentialsProvider(
		accessKeyID,
		secretAccessKey,
		"",
	)

	awsCfg, err := awsConfig.LoadDefaultConfig(
		ctx,
		awsConfig.WithRegion(region),
		awsConfig.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := s3.NewFromConfig(awsCfg)
	return &S3Client{client: client}, nil
}

// GetObject retrieves an object from S3
func (c *S3Client) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := c.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting object from S3: %w", err)
	}

	return result.Body, nil
}

func (c *S3Client) ListBuckets(ctx context.Context) ([]string, error) {
	input := &s3.ListBucketsInput{}
	result, err := c.client.ListBuckets(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error listing buckets: %w", err)
	}

	buckets := make([]string, len(result.Buckets))
	for i, bucket := range result.Buckets {
		buckets[i] = *bucket.Name
	}
	return buckets, nil
}

func (c *S3Client) CheckHealth(ctx context.Context) error {
	_, err := c.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("error listing buckets: %w", err)
	}
	return nil
}
