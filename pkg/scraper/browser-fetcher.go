package scraper

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

// BrowserFetcher fetches fully rendered HTML from a CDP-compatible browser.
// It connects to an already-running browser instance (e.g. Lightpanda or
// Chrome/Chromium started with --remote-debugging-port) and returns the outer
// HTML of the page after JavaScript execution.
type BrowserFetcher struct {
	endpoint string
}

// newBrowserFetcher creates a BrowserFetcher from its configuration.
func newBrowserFetcher(conf *BrowserConfiguration) *BrowserFetcher {
	return &BrowserFetcher{endpoint: conf.Endpoint}
}

// GetHTML connects to the browser, navigates to rawurl, and returns the fully
// rendered outer HTML. The provided context controls the overall deadline; the
// caller is responsible for setting an appropriate timeout.
//
// chromedp.NoModifyURL is required for Lightpanda (and harmless for Chrome):
// it disables chromedp's default behaviour of fetching /json/version and
// rewriting the WebSocket URL to a per-target path.
func (bf *BrowserFetcher) GetHTML(ctx context.Context, rawurl string) (string, error) {
	allocCtx, cancelAlloc := chromedp.NewRemoteAllocator(ctx, bf.endpoint, chromedp.NoModifyURL)
	defer cancelAlloc()

	tabCtx, cancelTab := chromedp.NewContext(allocCtx)
	defer cancelTab()

	var result string
	if err := chromedp.Run(tabCtx,
		chromedp.Navigate(rawurl),
		chromedp.OuterHTML("html", &result),
	); err != nil {
		return "", fmt.Errorf("browser fetch of %q failed: %w", rawurl, err)
	}
	return result, nil
}
