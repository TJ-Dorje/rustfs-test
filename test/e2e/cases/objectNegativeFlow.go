package cases

import (
	"context"
	"errors"
	"rustfs-e2e/configuration"
	"rustfs-e2e/helpers"
	"testing"

	"github.com/aws/smithy-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ValidateGetNonExistentObject(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()
	key := helpers.RandomKey()

	err := helpers.CreateBucket(ctx, client, bucket)
	require.NoError(t, err, "create bucket failed")

	_, err = helpers.GetObject(ctx, client, bucket, key)
	assert.Error(t, err, "getting non-existent object should fail")

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		t.Logf("error code (expected NoSuchKey): %s", apiErr.ErrorCode())
		assert.Equal(t, "NoSuchKey", apiErr.ErrorCode())
	}

	if cfg.TeardownEnabled {
		_ = helpers.DeleteBucket(ctx, client, bucket)
	}
}

func ValidateUploadToNonExistentBucket(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()
	key := helpers.RandomKey()

	err := helpers.PutObject(ctx, client, bucket, key, helpers.RandomPayload(256))
	assert.Error(t, err, "uploading to non-existent bucket should fail")

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		t.Logf("error code (expected NoSuchBucket): %s", apiErr.ErrorCode())
		assert.Equal(t, "NoSuchBucket", apiErr.ErrorCode())
	}
}

func ValidateGetFromNonExistentBucket(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()
	key := helpers.RandomKey()

	_, err := helpers.GetObject(ctx, client, bucket, key)
	assert.Error(t, err, "getting object from non-existent bucket should fail")
	t.Logf("get from non-existent bucket error (expected): %v", err)
}

func ValidateDeleteNonExistentObject(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()
	key := helpers.RandomKey()

	err := helpers.CreateBucket(ctx, client, bucket)
	require.NoError(t, err, "create bucket failed")

	// S3 spec: DeleteObject on non-existent key returns 204, not error
	err = helpers.DeleteObject(ctx, client, bucket, key)
	assert.NoError(t, err, "S3 spec: delete non-existent object should return 204, not error")
	t.Logf("delete non-existent object: no error (S3 spec compliant)")

	if cfg.TeardownEnabled {
		_ = helpers.DeleteBucket(ctx, client, bucket)
	}
}
