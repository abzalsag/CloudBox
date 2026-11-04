package s3

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"CloudBox/internal/entity"
)

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
}

type Storage struct {
	client *awss3.Client
	bucket string
}

func New(ctx context.Context, cfg Config) (*Storage, error) {
	client := awss3.New(awss3.Options{
		Region:                     "us-east-1",
		BaseEndpoint:               aws.String(cfg.Endpoint),
		Credentials:                credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		UsePathStyle:               true,
		RequestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
		ResponseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
	})

	s := &Storage{client: client, bucket: cfg.Bucket}
	if err := s.ensureBucket(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Storage) ensureBucket(ctx context.Context) error {
	_, err := s.client.HeadBucket(ctx, &awss3.HeadBucketInput{Bucket: aws.String(s.bucket)})
	if err == nil {
		return nil
	}

	_, err = s.client.CreateBucket(ctx, &awss3.CreateBucketInput{Bucket: aws.String(s.bucket)})
	if err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) || errors.As(err, &exists) {
			return nil
		}
		return fmt.Errorf("create bucket %q: %w", s.bucket, err)
	}
	return nil
}

func (s *Storage) Bucket() string { return s.bucket }

func (s *Storage) Put(ctx context.Context, key string, r io.ReadSeeker, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          r,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("s3 put %q: %w", key, err)
	}
	return nil
}

func (s *Storage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var noKey *types.NoSuchKey
		if errors.As(err, &noKey) {
			return nil, entity.ErrNotFound
		}
		return nil, fmt.Errorf("s3 get %q: %w", key, err)
	}
	return out.Body, nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &awss3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("s3 delete %q: %w", key, err)
	}
	return nil
}
