package scraper

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

// ScrapeAndClean uses headless Chrome to render JS-heavy pages, extracts the
// page HTML, prunes scripts/styles, and returns a cleaned text body.
func ScrapeAndClean(url string) (rawHTML string, bodyHTML string, cleaned string, err error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Add a timeout to avoid hanging
	ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	var html string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		// wait for network to be quiet-ish
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	}

	if err = chromedp.Run(ctx, tasks); err != nil {
		return "", "", "", err
	}

	// Parse with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html, "", "", err
	}

	bodySel := doc.Find("body")
	if bodySel.Length() == 0 {
		return html, "", "", errors.New("no <body> found")
	}

	// Clone body to manipulate
	tmp := goquery.NewDocumentFromNode(bodySel.Get(0))
	// Remove scripts and styles
	tmp.Find("script, style, noscript").Remove()

	bodyHTML, _ = tmp.Html()

	clean := strings.TrimSpace(tmp.Text())
	clean = strings.Join(strings.Fields(clean), " ")

	return html, bodyHTML, clean, nil
}