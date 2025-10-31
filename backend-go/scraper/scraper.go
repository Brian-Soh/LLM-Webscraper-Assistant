package scraper

import (
	"bytes"
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

	// Extract text with line breaks
	var buf bytes.Buffer
	tmp.Find("p, h1, h2, h3, h4, h5, h6, li, a, td, th, div, span").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			buf.WriteString(text)
			buf.WriteString("\n")
		}
	})
	clean := normalizeWhitespace(buf.String())

	// Get the cleaned body HTML (no scripts/styles)
	bodySel.Find("script, style, noscript").Remove()
	bodyHTML, _ = bodySel.Html()

	return html, bodyHTML, clean, nil
}

func normalizeWhitespace(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		// Collapse inner spaces
		ln = strings.Join(strings.Fields(ln), " ")
		out = append(out, ln)
	}
	return strings.Join(out, "\n")
}
