package configuration

import "github.com/mxschmitt/playwright-go"

var (
	pw      *playwright.Playwright
	browser playwright.Browser
)

func InitPlaywright(cfg *Config) error {
	var err error
	pw, err = playwright.Run()
	if err != nil {
		return err
	}
	opts := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(cfg.Headless),
	}
	if cfg.NoSandbox {
		opts.Args = []string{"--no-sandbox", "--disable-dev-shm-usage"}
	}
	browser, err = pw.Chromium.Launch(opts)
	return err
}

func StopPlaywright() {
	if browser != nil {
		browser.Close()
	}
	if pw != nil {
		pw.Stop()
	}
}

func NewPage() (playwright.Page, error) {
	return browser.NewPage()
}
