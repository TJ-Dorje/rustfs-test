package cases

import (
	"context"
	"rustfs-e2e/configuration"
	"rustfs-e2e/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ValidateBucketCreation(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()

	err := helpers.CreateBucket(ctx, client, bucket)
	require.NoError(t, err, "create bucket failed")

	exists, err := helpers.BucketExists(ctx, client, bucket)
	require.NoError(t, err)
	assert.True(t, exists, "bucket should exist after creation")

	if cfg.TeardownEnabled {
		_ = helpers.DeleteBucket(ctx, client, bucket)
	}
}

func ValidateBucketDeletion(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()

	err := helpers.CreateBucket(ctx, client, bucket)
	require.NoError(t, err, "create bucket failed")

	err = helpers.DeleteBucket(ctx, client, bucket)
	require.NoError(t, err, "delete bucket failed")

	exists, err := helpers.BucketExists(ctx, client, bucket)
	require.NoError(t, err)
	assert.False(t, exists, "bucket should not exist after deletion")
}

func ValidateBucketLifecycle(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()

	err := helpers.CreateBucket(ctx, client, bucket)
	require.NoError(t, err, "create bucket failed")
	t.Logf("created bucket: %s", bucket)

	exists, err := helpers.BucketExists(ctx, client, bucket)
	require.NoError(t, err)
	assert.True(t, exists, "bucket should exist after creation")

	if cfg.TeardownEnabled {
		err = helpers.DeleteBucket(ctx, client, bucket)
		require.NoError(t, err, "delete bucket failed")
		t.Logf("deleted bucket: %s", bucket)

		exists, err = helpers.BucketExists(ctx, client, bucket)
		require.NoError(t, err)
		assert.False(t, exists, "bucket should not exist after deletion")
	}
}

func ValidateListBuckets(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()

	err := helpers.CreateBucket(ctx, client, bucket)
	require.NoError(t, err, "create bucket failed")

	resp, err := client.ListBuckets(ctx, nil)
	require.NoError(t, err, "list buckets failed")

	names := make([]string, 0, len(resp.Buckets))
	for _, b := range resp.Buckets {
		if b.Name != nil {
			names = append(names, *b.Name)
		}
	}
	assert.Contains(t, names, bucket, "created bucket should appear in list")
	t.Logf("buckets found: %v", names)

	if cfg.TeardownEnabled {
		_ = helpers.DeleteBucket(ctx, client, bucket)
	}
}
