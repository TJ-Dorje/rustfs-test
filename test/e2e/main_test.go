package main_test

import (
	"context"
	"log"
	"os"
	"rustfs-e2e/cases"
	"rustfs-e2e/configuration"
	"testing"
)

var (
	ctx       context.Context
	AppConfig *configuration.Config
)

func init() {
	ctx = context.Background()
	if err := configuration.LoadEnvFile(configuration.EnvFile); err != nil {
		log.Fatalf("error loading env file: %v", err)
	}
	AppConfig = configuration.AppConfig()
}

func TestMain(m *testing.M) {
	log.Println("starting rustfs e2e test suite")
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestBucketHappyFlowSuite(t *testing.T) {
	t.Run("BucketCreation", cases.ValidateBucketCreation)
	t.Run("BucketDeletion", cases.ValidateBucketDeletion)
	t.Run("BucketLifecycle", cases.ValidateBucketLifecycle)
	t.Run("ListBuckets", cases.ValidateListBuckets)
}

func TestBucketNegativeFlowSuite(t *testing.T) {
	t.Run("DuplicateBucketCreation", cases.ValidateDuplicateBucketCreation)
	t.Run("DeleteNonEmptyBucket", cases.ValidateDeleteNonEmptyBucket)
	t.Run("DeleteNonExistentBucket", cases.ValidateDeleteNonExistentBucket)
	t.Run("HeadNonExistentBucket", cases.ValidateHeadNonExistentBucket)
}

func TestObjectHappyFlowSuite(t *testing.T) {
	t.Run("ObjectUpload", cases.ValidateObjectUpload)
	t.Run("ObjectDownload", cases.ValidateObjectDownload)
	t.Run("ObjectDeletion", cases.ValidateObjectDeletion)
	t.Run("ObjectLifecycle", cases.ValidateObjectLifecycle)
	t.Run("ListObjects", cases.ValidateListObjects)
}

func TestObjectNegativeFlowSuite(t *testing.T) {
	t.Run("GetNonExistentObject", cases.ValidateGetNonExistentObject)
	t.Run("UploadToNonExistentBucket", cases.ValidateUploadToNonExistentBucket)
	t.Run("GetFromNonExistentBucket", cases.ValidateGetFromNonExistentBucket)
	t.Run("DeleteNonExistentObject", cases.ValidateDeleteNonExistentObject)
}
