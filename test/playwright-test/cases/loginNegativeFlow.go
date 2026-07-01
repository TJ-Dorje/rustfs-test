package cases

import (
	"rustfs-pw/configuration"

	"rustfs-pw/helpers/actions"
	"rustfs-pw/helpers/components"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ValidateLoginInvalidCredentials(t *testing.T) {
	cfg := configuration.AppConfig()
	page, err := configuration.NewPage()
	require.NoError(t, err, "create page failed")
	defer page.Close()

	badCfg := *cfg
	badCfg.Password = "wrongpassword"

	err = actions.LoginFill(page, &badCfg)
	require.NoError(t, err, "login action itself failed unexpectedly")

	url := page.URL()
	assert.Contains(t, url, "login", "should stay on login page after bad credentials")
	t.Logf("url after bad login: %s", url)

	errVisible, _ := components.ToastVisible(page, components.ToastError, "Login Failed")
	assert.True(t, errVisible, "login failed toast should be visible after bad credentials")
}
