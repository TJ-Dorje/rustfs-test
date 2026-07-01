package configuration

import "os"

type Config struct {
	Initialized     bool
	ConsoleURL      string
	Username        string
	Password        string
	Headless        bool
	TeardownEnabled bool
	NoSandbox       bool
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
	appConfig.ConsoleURL      = getEnv("RUSTFS_CONSOLE_URL", DefaultConsoleURL)
	appConfig.Username        = getEnv("RUSTFS_USERNAME", DefaultUsername)
	appConfig.Password        = getEnv("RUSTFS_PASSWORD", DefaultPassword)
	appConfig.Headless        = getEnv("HEADLESS", "true") != "false"
	appConfig.TeardownEnabled = getEnv("TEARDOWN_ENABLED", "true") != "false"
	appConfig.NoSandbox       = getEnv("NO_SANDBOX", "false") == "true"
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
