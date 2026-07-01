package configuration

import "time"

const (
	DefaultConsoleURL = "http://localhost:9001"
	DefaultConsolePath = "/rustfs/console"
	LoginPath          = "/rustfs/console/auth/login"
	BrowserPath        = "/rustfs/console/browser"
	UsersPath          = "/rustfs/console/users/"
	UserGroupsPath     = "/rustfs/console/user-groups/"
	DefaultUsername    = "rustfsadmin"
	DefaultPassword    = "rustfsadmin"
	DefaultTimeout     = 30 * time.Second
	EnvFile            = ".env"
)
