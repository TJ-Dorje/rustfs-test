package actions

import (
	"fmt"
	"rustfs-pw/configuration"
	"rustfs-pw/helpers/components"

	"github.com/mxschmitt/playwright-go"
)

func NavigateToUsers(page playwright.Page, cfg *configuration.Config) error {
	if _, err := page.Goto(cfg.ConsoleURL + configuration.UsersPath); err != nil {
		return fmt.Errorf("navigate to users: %w", err)
	}
	return page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})
}

func CreateUserViaUI(page playwright.Page, username, password string) error {
	if err := page.Locator("button:has-text(\"Add User\")").Click(); err != nil {
		return fmt.Errorf("click add user: %w", err)
	}
	if err := page.Locator("[role=\"dialog\"] input").Nth(0).WaitFor(); err != nil {
		return fmt.Errorf("wait for username input: %w", err)
	}
	if err := page.Locator("[role=\"dialog\"] input").Nth(0).Fill(username); err != nil {
		return fmt.Errorf("fill username: %w", err)
	}
	if err := page.Locator("[role=\"dialog\"] input").Nth(1).Fill(password); err != nil {
		return fmt.Errorf("fill password: %w", err)
	}
	if err := page.Locator("[role=\"dialog\"] button:text-is(\"Submit\")").Click(); err != nil {
		return fmt.Errorf("click submit: %w", err)
	}
	return page.Locator("[data-slot=\"dialog-content\"]").WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateDetached,
	})
}

func DeleteUserViaUI(page playwright.Page, username string) error {
	if err := components.ClickRowAction(page, username, "Delete"); err != nil {
		return fmt.Errorf("click delete for user %s: %w", username, err)
	}
	if err := components.ConfirmAlertDialog(page); err != nil {
		return fmt.Errorf("confirm delete dialog: %w", err)
	}
	return components.WaitForRowDetached(page, username)
}

func UserExistsInUI(page playwright.Page, username string) (bool, error) {
	return components.RowExists(page, username)
}

func UserMemberOfGroup(page playwright.Page, username, groupName string) (bool, error) {
	selector := fmt.Sprintf("tr:has-text(\"%s\"):has-text(\"%s\")", username, groupName)
	count, err := page.Locator(selector).Count()
	return count > 0, err
}
