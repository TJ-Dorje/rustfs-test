package actions

import (
	"fmt"
	"rustfs-pw/configuration"
	"rustfs-pw/helpers/components"

	"github.com/mxschmitt/playwright-go"
)

func NavigateToBuckets(page playwright.Page, cfg *configuration.Config) error {
	if _, err := page.Goto(cfg.ConsoleURL + configuration.BrowserPath); err != nil {
		return fmt.Errorf("navigate to buckets: %w", err)
	}
	return page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})
}

func CreateBucketViaUI(page playwright.Page, name string) error {
	if err := page.Locator("button:has-text(\"Create Bucket\")").Click(); err != nil {
		return fmt.Errorf("click create bucket: %w", err)
	}
	if err := page.Locator("#bucket-name").WaitFor(); err != nil {
		return fmt.Errorf("wait for bucket name input: %w", err)
	}
	if err := page.Locator("#bucket-name").Fill(name); err != nil {
		return fmt.Errorf("fill bucket name: %w", err)
	}
	if err := page.Locator("button:text-is(\"Create\")").Click(); err != nil {
		return fmt.Errorf("click create: %w", err)
	}
	return components.WaitForToast(page, components.ToastSuccess, "Create Success")
}

func DeleteBucketViaUI(page playwright.Page, name string) error {
	if err := components.ClickRowAction(page, name, "Delete"); err != nil {
		return fmt.Errorf("click delete for bucket %s: %w", name, err)
	}
	if err := components.ConfirmAlertDialog(page); err != nil {
		return fmt.Errorf("confirm delete dialog: %w", err)
	}
	return components.WaitForRowDetached(page, name)
}

func BucketExistsInUI(page playwright.Page, name string) (bool, error) {
	return components.RowExists(page, name)
}
