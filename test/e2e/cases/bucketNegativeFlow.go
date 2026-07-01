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

func ValidateDuplicateBucketCreation(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()

	err := helpers.CreateBucket(ctx, client, bucket)
	require.NoError(t, err, "first create bucket failed")

	err = helpers.CreateBucket(ctx, client, bucket)
	// RustFS (us-east-1): duplicate create returns 200, no error — S3 spec allows this for same owner
	// other regions/implementations may return BucketAlreadyExists or BucketAlreadyOwnedByYou
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			validCodes := []string{"BucketAlreadyExists", "BucketAlreadyOwnedByYou"}
			t.Logf("duplicate bucket error code: %s", apiErr.ErrorCode())
			assert.Contains(t, validCodes, apiErr.ErrorCode())
		}
	} else {
		t.Logf("duplicate bucket creation: no error (RustFS us-east-1 behavior)")
	}

	if cfg.TeardownEnabled {
		_ = helpers.DeleteBucket(ctx, client, bucket)
	}
}

func ValidateDeleteNonEmptyBucket(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()
	key := helpers.RandomKey()

	err := helpers.CreateBucket(ctx, client, bucket)
	require.NoError(t, err, "create bucket failed")

	err = helpers.PutObject(ctx, client, bucket, key, helpers.RandomPayload(256))
	require.NoError(t, err, "put object failed")

	err = helpers.DeleteBucket(ctx, client, bucket)
	assert.Error(t, err, "deleting non-empty bucket should fail")

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		t.Logf("error code (expected BucketNotEmpty): %s", apiErr.ErrorCode())
		assert.Equal(t, "BucketNotEmpty", apiErr.ErrorCode())
	}

	if cfg.TeardownEnabled {
		_ = helpers.DeleteObject(ctx, client, bucket, key)
		_ = helpers.DeleteBucket(ctx, client, bucket)
	}
}

func ValidateDeleteNonExistentBucket(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()

	err := helpers.DeleteBucket(ctx, client, bucket)
	assert.Error(t, err, "deleting non-existent bucket should fail")
	t.Logf("non-existent bucket delete error (expected): %v", err)
}

func ValidateHeadNonExistentBucket(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()

	exists, err := helpers.BucketExists(ctx, client, bucket)
	require.NoError(t, err)
	assert.False(t, exists, "non-existent bucket should not be found")
	t.Logf("head non-existent bucket: exists=%v", exists)
}
