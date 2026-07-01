package components

import (
	"fmt"

	"github.com/mxschmitt/playwright-go"
)

// RowExists checks if a table row containing the given text is present
func RowExists(page playwright.Page, text string) (bool, error) {
	locator := page.Locator(fmt.Sprintf("tr:has-text(\"%s\")", text))
	count, err := locator.Count()
	return count > 0, err
}

// WaitForRowDetached waits until a table row containing the given text is removed from DOM
func WaitForRowDetached(page playwright.Page, text string) error {
	_, err := page.WaitForSelector(
		fmt.Sprintf("tr:has-text(\"%s\")", text),
		playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached},
	)
	return err
}

// ClickRowAction clicks an action button within a specific table row
func ClickRowAction(page playwright.Page, rowText, buttonText string) error {
	selector := fmt.Sprintf("tr:has-text(\"%s\") button:has-text(\"%s\")", rowText, buttonText)
	if err := page.Click(selector); err != nil {
		return fmt.Errorf("click %q in row %q: %w", buttonText, rowText, err)
	}
	return nil
}
