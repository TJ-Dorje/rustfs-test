package configuration

import "os"

type Config struct {
	Initialized     bool
	Endpoint        string
	Region          string
	AccessKey       string
	SecretKey       string
	TeardownEnabled bool
}

var appConfig Config

func AppConfig() *Config {
	if !appConfig.Initialized {
		loadEnvironmentVariables()
		appConfig.Initialized = true
	}
	return &appConfig
}

func loadEnvironmentVariables() {
	appConfig.Endpoint = getEnv("RUSTFS_ENDPOINT", DefaultEndpoint)
	appConfig.Region = getEnv("AWS_DEFAULT_REGION", DefaultRegion)
	appConfig.AccessKey = getEnv("AWS_ACCESS_KEY_ID", "rustfsadmin")
	appConfig.SecretKey = getEnv("AWS_SECRET_ACCESS_KEY", "rustfsadmin")
	appConfig.TeardownEnabled = getEnv("TEARDOWN_ENABLED", "true") != "false"
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
