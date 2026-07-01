package cases

import (
	"context"
	"rustfs-e2e/configuration"
	"rustfs-e2e/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ValidateObjectUpload(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()
	key := helpers.RandomKey()
	payload := helpers.RandomPayload(1024)

	err := helpers.CreateBucket(ctx, client, bucket)
	require.NoError(t, err, "create bucket failed")

	err = helpers.PutObject(ctx, client, bucket, key, payload)
	require.NoError(t, err, "put object failed")

	exists, err := helpers.ObjectExists(ctx, client, bucket, key)
	require.NoError(t, err)
	assert.True(t, exists, "object should exist after upload")
	t.Logf("uploaded object: %s/%s", bucket, key)

	if cfg.TeardownEnabled {
		_ = helpers.DeleteObject(ctx, client, bucket, key)
		_ = helpers.DeleteBucket(ctx, client, bucket)
	}
}

func ValidateObjectDownload(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()
	key := helpers.RandomKey()
	payload := helpers.RandomPayload(1024)

	err := helpers.CreateBucket(ctx, client, bucket)
	require.NoError(t, err, "create bucket failed")

	err = helpers.PutObject(ctx, client, bucket, key, payload)
	require.NoError(t, err, "put object failed")

	downloaded, err := helpers.GetObject(ctx, client, bucket, key)
	require.NoError(t, err, "get object failed")

	assert.Equal(t, payload, downloaded, "downloaded content should match uploaded content")
	t.Logf("verified object content: %s/%s (%d bytes)", bucket, key, len(downloaded))

	if cfg.TeardownEnabled {
		_ = helpers.DeleteObject(ctx, client, bucket, key)
		_ = helpers.DeleteBucket(ctx, client, bucket)
	}
}

func ValidateObjectDeletion(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()
	key := helpers.RandomKey()
	payload := helpers.RandomPayload(1024)

	err := helpers.CreateBucket(ctx, client, bucket)
	require.NoError(t, err, "create bucket failed")

	err = helpers.PutObject(ctx, client, bucket, key, payload)
	require.NoError(t, err, "put object failed")

	err = helpers.DeleteObject(ctx, client, bucket, key)
	require.NoError(t, err, "delete object failed")

	exists, err := helpers.ObjectExists(ctx, client, bucket, key)
	require.NoError(t, err)
	assert.False(t, exists, "object should not exist after deletion")

	if cfg.TeardownEnabled {
		_ = helpers.DeleteBucket(ctx, client, bucket)
	}
}

func ValidateObjectLifecycle(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()
	key := helpers.RandomKey()
	payload := helpers.RandomPayload(1024)

	err := helpers.CreateBucket(ctx, client, bucket)
	require.NoError(t, err, "create bucket failed")

	err = helpers.PutObject(ctx, client, bucket, key, payload)
	require.NoError(t, err, "put object failed")
	t.Logf("uploaded object: %s/%s", bucket, key)

	downloaded, err := helpers.GetObject(ctx, client, bucket, key)
	require.NoError(t, err, "get object failed")
	assert.Equal(t, payload, downloaded, "downloaded content should match uploaded content")

	if cfg.TeardownEnabled {
		err = helpers.DeleteObject(ctx, client, bucket, key)
		require.NoError(t, err, "delete object failed")

		exists, err := helpers.ObjectExists(ctx, client, bucket, key)
		require.NoError(t, err)
		assert.False(t, exists, "object should not exist after deletion")
		t.Logf("deleted object: %s/%s", bucket, key)

		_ = helpers.DeleteBucket(ctx, client, bucket)
	}
}

func ValidateListObjects(t *testing.T) {
	cfg := configuration.AppConfig()
	client := helpers.NewS3Client(cfg)
	ctx := context.Background()

	bucket := helpers.RandomBucketName()
	payload := helpers.RandomPayload(256)

	err := helpers.CreateBucket(ctx, client, bucket)
	require.NoError(t, err, "create bucket failed")

	keys := []string{helpers.RandomKey(), helpers.RandomKey(), helpers.RandomKey()}
	for _, k := range keys {
		err = helpers.PutObject(ctx, client, bucket, k, payload)
		require.NoError(t, err, "put object failed for key: %s", k)
	}

	listed, err := helpers.ListObjects(ctx, client, bucket)
	require.NoError(t, err, "list objects failed")

	for _, k := range keys {
		assert.Contains(t, listed, k, "key should appear in list: %s", k)
	}
	t.Logf("listed %d objects in bucket %s", len(listed), bucket)

	if cfg.TeardownEnabled {
		for _, k := range keys {
			_ = helpers.DeleteObject(ctx, client, bucket, k)
		}
		_ = helpers.DeleteBucket(ctx, client, bucket)
	}
}
