package actions

import (
	"fmt"
	"rustfs-pw/configuration"
	"rustfs-pw/helpers/components"

	"github.com/mxschmitt/playwright-go"
)

func NavigateToGroups(page playwright.Page, cfg *configuration.Config) error {
	if _, err := page.Goto(cfg.ConsoleURL + configuration.UserGroupsPath); err != nil {
		return fmt.Errorf("navigate to groups: %w", err)
	}
	return page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})
}

// CreateGroupViaUI creates a user group and optionally adds members.
// Members are selected via the cmdk dropdown which requires DispatchEvent to register clicks.
func CreateGroupViaUI(page playwright.Page, groupName string, members []string) error {
	if err := page.Locator("button:has-text(\"Add User Group\")").Click(); err != nil {
		return fmt.Errorf("click add user group: %w", err)
	}
	if err := page.Locator("[role=\"dialog\"] input").First().WaitFor(); err != nil {
		return fmt.Errorf("wait for name input: %w", err)
	}
	if err := page.Locator("[role=\"dialog\"] input").First().Fill(groupName); err != nil {
		return fmt.Errorf("fill group name: %w", err)
	}
	for _, member := range members {
		if err := page.Locator("[role=\"dialog\"] button:has-text(\"Select user group members\")").Click(); err != nil {
			return fmt.Errorf("open users dropdown: %w", err)
		}
		if err := page.Locator("[role=\"listbox\"]").WaitFor(); err != nil {
			return fmt.Errorf("wait for user listbox: %w", err)
		}
		selector := fmt.Sprintf("[role=\"option\"]:has-text(\"%s\")", member)
		if err := page.Locator(selector).DispatchEvent("click", nil); err != nil {
			return fmt.Errorf("select member %s: %w", member, err)
		}
	}
	if err := page.Locator("[role=\"dialog\"] button:text-is(\"Submit\")").Click(); err != nil {
		return fmt.Errorf("click submit: %w", err)
	}
	return page.Locator("[data-slot=\"dialog-content\"]").WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateDetached,
	})
}

func DeleteGroupViaUI(page playwright.Page, groupName string) error {
	if err := components.ClickRowAction(page, groupName, "Delete"); err != nil {
		return fmt.Errorf("click delete for group %s: %w", groupName, err)
	}
	if err := components.DeleteAlertDialog(page); err != nil {
		return fmt.Errorf("confirm delete dialog: %w", err)
	}
	return components.WaitForRowDetached(page, groupName)
}

func GroupExistsInUI(page playwright.Page, groupName string) (bool, error) {
	return components.RowExists(page, groupName)
}
