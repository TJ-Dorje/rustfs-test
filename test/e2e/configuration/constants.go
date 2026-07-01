package configuration

import "time"

const (
	DefaultEndpoint = "http://localhost:9000"
	DefaultRegion   = "us-east-1"
	DefaultTimeout  = 30 * time.Second
	EnvFile         = ".env"
)
