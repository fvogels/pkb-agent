package extern

import "github.com/pkg/browser"

func OpenURLInBrowser(url string) error {
	return browser.OpenURL(url)
}
