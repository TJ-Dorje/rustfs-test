package components

import (
	"fmt"

	"github.com/mxschmitt/playwright-go"
)

const (
	alertDialog        = `[role="alertdialog"]`
	alertDialogTitle   = `[data-slot="alert-dialog-title"]`
	alertDialogConfirm = `[role="alertdialog"] button:has-text("Confirm")`
	alertDialogCancel  = `[role="alertdialog"] button:has-text("Cancel")`
)

func AlertDialogVisible(page playwright.Page) (bool, error) {
	return page.Locator(alertDialog).IsVisible()
}

func AlertDialogTitle(page playwright.Page) (string, error) {
	return page.Locator(alertDialogTitle).InnerText()
}

func ConfirmAlertDialog(page playwright.Page) error {
	if _, err := page.WaitForSelector(alertDialogConfirm); err != nil {
		return fmt.Errorf("confirm button not found: %w", err)
	}
	return page.Click(alertDialogConfirm)
}

func CancelAlertDialog(page playwright.Page) error {
	if _, err := page.WaitForSelector(alertDialogCancel); err != nil {
		return fmt.Errorf("cancel button not found: %w", err)
	}
	return page.Click(alertDialogCancel)
}

// DeleteAlertDialog clicks the "Delete" button in an alertdialog.
// Some confirm dialogs use "Delete" instead of "Confirm" (e.g. group deletion).
func DeleteAlertDialog(page playwright.Page) error {
	selector := `[role="alertdialog"] button:has-text("Delete")`
	if _, err := page.WaitForSelector(selector); err != nil {
		return fmt.Errorf("delete button not found: %w", err)
	}
	return page.Click(selector)
}
