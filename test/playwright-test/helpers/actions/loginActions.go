package actions

import (
	"fmt"
	"rustfs-pw/configuration"

	"github.com/mxschmitt/playwright-go"
)

func Login(page playwright.Page, cfg *configuration.Config) error {
	if _, err := page.Goto(cfg.ConsoleURL + configuration.LoginPath); err != nil {
		return fmt.Errorf("navigate to login: %w", err)
	}
	if err := page.Locator("#accessKey").WaitFor(); err != nil {
		return fmt.Errorf("wait for login form: %w", err)
	}
	if err := page.Locator("#accessKey").Fill(cfg.Username); err != nil {
		return fmt.Errorf("fill username: %w", err)
	}
	if err := page.Locator("#secretKey").Fill(cfg.Password); err != nil {
		return fmt.Errorf("fill password: %w", err)
	}
	if err := page.Locator("button[type=submit]").Click(); err != nil {
		return fmt.Errorf("click login: %w", err)
	}
	if err := page.WaitForURL("**/browser/**"); err != nil {
		return fmt.Errorf("login redirect failed: %w", err)
	}
	return nil
}

// LoginFill fills credentials and submits without waiting for redirect.
// Use in negative cases where login is expected to fail.
func LoginFill(page playwright.Page, cfg *configuration.Config) error {
	if _, err := page.Goto(cfg.ConsoleURL + configuration.LoginPath); err != nil {
		return fmt.Errorf("navigate to login: %w", err)
	}
	if err := page.Locator("#accessKey").WaitFor(); err != nil {
		return fmt.Errorf("wait for login form: %w", err)
	}
	if err := page.Locator("#accessKey").Fill(cfg.Username); err != nil {
		return fmt.Errorf("fill username: %w", err)
	}
	if err := page.Locator("#secretKey").Fill(cfg.Password); err != nil {
		return fmt.Errorf("fill password: %w", err)
	}
	if err := page.Locator("button[type=submit]").Click(); err != nil {
		return fmt.Errorf("click login: %w", err)
	}
	page.WaitForTimeout(2000)
	return nil
}
