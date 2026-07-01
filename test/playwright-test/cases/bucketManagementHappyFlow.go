package cases

import (
	"fmt"
	"rustfs-pw/configuration"

	"rustfs-pw/helpers/actions"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ValidateBucketCreationViaUI(t *testing.T) {
	cfg := configuration.AppConfig()
	page, err := configuration.NewPage()
	require.NoError(t, err, "create page failed")
	defer page.Close()

	err = actions.Login(page, cfg)
	require.NoError(t, err, "login failed")

	bucket := fmt.Sprintf("ui-bucket-%d", time.Now().Unix())

	err = actions.CreateBucketViaUI(page, bucket)
	require.NoError(t, err, "create bucket via UI failed")
	t.Logf("created bucket via UI: %s", bucket)

	err = actions.NavigateToBuckets(page, cfg)
	require.NoError(t, err)

	exists, err := actions.BucketExistsInUI(page, bucket)
	require.NoError(t, err)
	assert.True(t, exists, "bucket should appear in list after creation")

	if cfg.TeardownEnabled {
		err = actions.DeleteBucketViaUI(page, bucket)
		require.NoError(t, err, "delete bucket via UI failed")
		t.Logf("deleted bucket via UI: %s", bucket)
	}
}

func ValidateBucketDeletionViaUI(t *testing.T) {
	cfg := configuration.AppConfig()
	page, err := configuration.NewPage()
	require.NoError(t, err, "create page failed")
	defer page.Close()

	err = actions.Login(page, cfg)
	require.NoError(t, err, "login failed")

	bucket := fmt.Sprintf("ui-bucket-%d", time.Now().Unix())

	err = actions.CreateBucketViaUI(page, bucket)
	require.NoError(t, err, "create bucket failed")

	err = actions.NavigateToBuckets(page, cfg)
	require.NoError(t, err)

	err = actions.DeleteBucketViaUI(page, bucket)
	require.NoError(t, err, "delete bucket via UI failed")
	t.Logf("deleted bucket via UI: %s", bucket)

	err = actions.NavigateToBuckets(page, cfg)
	require.NoError(t, err)

	exists, err := actions.BucketExistsInUI(page, bucket)
	require.NoError(t, err)
	assert.False(t, exists, "bucket should not appear in list after deletion")
}

func ValidateBucketLifecycleViaUI(t *testing.T) {
	cfg := configuration.AppConfig()
	page, err := configuration.NewPage()
	require.NoError(t, err, "create page failed")
	defer page.Close()

	err = actions.Login(page, cfg)
	require.NoError(t, err, "login failed")

	bucket := fmt.Sprintf("ui-bucket-%d", time.Now().Unix())

	err = actions.CreateBucketViaUI(page, bucket)
	require.NoError(t, err, "create bucket failed")
	t.Logf("created bucket via UI: %s", bucket)

	err = actions.NavigateToBuckets(page, cfg)
	require.NoError(t, err)

	exists, err := actions.BucketExistsInUI(page, bucket)
	require.NoError(t, err)
	assert.True(t, exists, "bucket should exist after creation")

	if cfg.TeardownEnabled {
		err = actions.DeleteBucketViaUI(page, bucket)
		require.NoError(t, err, "delete bucket failed")
		t.Logf("deleted bucket via UI: %s", bucket)

		err = actions.NavigateToBuckets(page, cfg)
		require.NoError(t, err)

		exists, err = actions.BucketExistsInUI(page, bucket)
		require.NoError(t, err)
		assert.False(t, exists, "bucket should not exist after deletion")
	}
}
