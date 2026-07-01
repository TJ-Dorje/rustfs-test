package components

import "github.com/mxschmitt/playwright-go"

const (
	ToastSuccess = "success"
	ToastWarning = "warning"
	ToastError   = "error"
)

func toastSelector(toastType, text string) string {
	if text != "" {
		return `.cn-toast[data-type="` + toastType + `"] [data-title]:has-text("` + text + `")`
	}
	return `.cn-toast[data-type="` + toastType + `"]`
}

func ToastVisible(page playwright.Page, toastType, text string) (bool, error) {
	return page.Locator(toastSelector(toastType, text)).IsVisible()
}

func WaitForToast(page playwright.Page, toastType, text string) error {
	_, err := page.WaitForSelector(toastSelector(toastType, text))
	return err
}
