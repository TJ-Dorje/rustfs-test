package helpers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"rustfs-e2e/configuration"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

func NewS3Client(cfg *configuration.Config) *s3.Client {
	c, _ := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	)
	return s3.NewFromConfig(c, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = true
	})
}

func CreateBucket(ctx context.Context, client *s3.Client, name string) error {
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	return err
}

func DeleteBucket(ctx context.Context, client *s3.Client, name string) error {
	_, err := client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	})
	return err
}

func BucketExists(ctx context.Context, client *s3.Client, name string) (bool, error) {
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NotFound" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func PutObject(ctx context.Context, client *s3.Client, bucket, key string, body []byte) error {
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})
	return err
}

func GetObject(ctx context.Context, client *s3.Client, bucket, key string) ([]byte, error) {
	resp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func DeleteObject(ctx context.Context, client *s3.Client, bucket, key string) error {
	_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

func ObjectExists(ctx context.Context, client *s3.Client, bucket, key string) (bool, error) {
	_, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NotFound" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ListObjects(ctx context.Context, client *s3.Client, bucket string) ([]string, error) {
	resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(resp.Contents))
	for _, obj := range resp.Contents {
		keys = append(keys, aws.ToString(obj.Key))
	}
	return keys, nil
}
