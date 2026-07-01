package main_test

import (
	"log"
	"os"
	"rustfs-pw/cases"
	"rustfs-pw/configuration"
	"testing"
)

var AppConfig *configuration.Config

func init() {
	if err := configuration.LoadEnvFile(configuration.EnvFile); err != nil {
		log.Fatalf("error loading env file: %v", err)
	}
	AppConfig = configuration.AppConfig()
	if err := configuration.InitPlaywright(AppConfig); err != nil {
		log.Fatalf("error initializing playwright: %v", err)
	}
}

func TestMain(m *testing.M) {
	log.Println("starting rustfs playwright gui test suite")
	exitVal := m.Run()
	configuration.StopPlaywright()
	os.Exit(exitVal)
}

func TestLoginHappyFlowSuite(t *testing.T) {
	t.Run("Login", cases.ValidateLogin)
}

func TestLoginNegativeFlowSuite(t *testing.T) {
	t.Run("InvalidCredentials", cases.ValidateLoginInvalidCredentials)
}

func TestBucketManagementHappyFlowSuite(t *testing.T) {
	t.Run("BucketCreation", cases.ValidateBucketCreationViaUI)
	t.Run("BucketDeletion", cases.ValidateBucketDeletionViaUI)
	t.Run("BucketLifecycle", cases.ValidateBucketLifecycleViaUI)
}

func TestUserManagementHappyFlowSuite(t *testing.T) {
	t.Run("UserCreation", cases.ValidateUserCreationViaUI)
	t.Run("GroupCreation", cases.ValidateGroupCreationViaUI)
	t.Run("AddUserToGroup", cases.ValidateAddUserToGroupViaUI)
}
