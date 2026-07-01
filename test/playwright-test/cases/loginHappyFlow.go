package cases

import (
	"rustfs-pw/configuration"

	"rustfs-pw/helpers/actions"
	"rustfs-pw/helpers/components"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ValidateLogin(t *testing.T) {
	cfg := configuration.AppConfig()
	page, err := configuration.NewPage()
	require.NoError(t, err, "create page failed")
	defer page.Close()

	err = actions.Login(page, cfg)
	require.NoError(t, err, "login failed")

	url := page.URL()
	assert.Contains(t, url, configuration.BrowserPath, "should land on browser after login")
	t.Logf("landed on: %s", url)

	visible, err := components.ToastVisible(page, components.ToastSuccess, "Login Success")
	require.NoError(t, err)
	assert.True(t, visible, "login success toast should be visible")
}
