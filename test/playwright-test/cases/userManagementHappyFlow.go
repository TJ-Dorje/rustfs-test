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

func ValidateUserCreationViaUI(t *testing.T) {
	cfg := configuration.AppConfig()
	page, err := configuration.NewPage()
	require.NoError(t, err, "create page failed")
	defer page.Close()

	err = actions.Login(page, cfg)
	require.NoError(t, err, "login failed")

	username := fmt.Sprintf("ui-user-%d", time.Now().Unix())

	err = actions.NavigateToUsers(page, cfg)
	require.NoError(t, err)

	err = actions.CreateUserViaUI(page, username, "Testpass123!")
	require.NoError(t, err, "create user via UI failed")
	t.Logf("created user via UI: %s", username)

	exists, err := actions.UserExistsInUI(page, username)
	require.NoError(t, err)
	assert.True(t, exists, "user should appear in list after creation")

	if cfg.TeardownEnabled {
		err = actions.DeleteUserViaUI(page, username)
		require.NoError(t, err, "delete user via UI failed")
		t.Logf("deleted user via UI: %s", username)
	}
}

func ValidateGroupCreationViaUI(t *testing.T) {
	cfg := configuration.AppConfig()
	page, err := configuration.NewPage()
	require.NoError(t, err, "create page failed")
	defer page.Close()

	err = actions.Login(page, cfg)
	require.NoError(t, err, "login failed")

	groupName := fmt.Sprintf("ui-group-%d", time.Now().Unix())

	err = actions.NavigateToGroups(page, cfg)
	require.NoError(t, err)

	err = actions.CreateGroupViaUI(page, groupName, nil)
	require.NoError(t, err, "create group via UI failed")
	t.Logf("created group via UI: %s", groupName)

	exists, err := actions.GroupExistsInUI(page, groupName)
	require.NoError(t, err)
	assert.True(t, exists, "group should appear in list after creation")

	if cfg.TeardownEnabled {
		err = actions.DeleteGroupViaUI(page, groupName)
		require.NoError(t, err, "delete group via UI failed")
		t.Logf("deleted group via UI: %s", groupName)
	}
}

func ValidateAddUserToGroupViaUI(t *testing.T) {
	cfg := configuration.AppConfig()
	page, err := configuration.NewPage()
	require.NoError(t, err, "create page failed")
	defer page.Close()

	err = actions.Login(page, cfg)
	require.NoError(t, err, "login failed")

	ts := time.Now().Unix()
	username := fmt.Sprintf("ui-user-%d", ts)
	groupName := fmt.Sprintf("ui-group-%d", ts)

	err = actions.NavigateToUsers(page, cfg)
	require.NoError(t, err)

	err = actions.CreateUserViaUI(page, username, "Testpass123!")
	require.NoError(t, err, "create user failed")
	t.Logf("created user: %s", username)

	err = actions.NavigateToGroups(page, cfg)
	require.NoError(t, err)

	err = actions.CreateGroupViaUI(page, groupName, []string{username})
	require.NoError(t, err, "create group with member failed")
	t.Logf("created group: %s with member: %s", groupName, username)

	err = actions.NavigateToUsers(page, cfg)
	require.NoError(t, err)

	isMember, err := actions.UserMemberOfGroup(page, username, groupName)
	require.NoError(t, err)
	assert.True(t, isMember, "user should be member of group")

	if cfg.TeardownEnabled {
		// Delete user first — removes group membership, leaving group empty
		err = actions.NavigateToUsers(page, cfg)
		require.NoError(t, err)
		err = actions.DeleteUserViaUI(page, username)
		require.NoError(t, err, "delete user failed")
		t.Logf("deleted user: %s", username)

		// Now delete the empty group
		err = actions.NavigateToGroups(page, cfg)
		require.NoError(t, err)
		err = actions.DeleteGroupViaUI(page, groupName)
		require.NoError(t, err, "delete group failed")
		t.Logf("deleted group: %s", groupName)
	}
}
