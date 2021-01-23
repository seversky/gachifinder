package scrape

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/seversky/gachifinder"
)

var _ gachifinder.Scraper = &PortalNews{}

// PortalNews struct.
type PortalNews struct {
	Scrape
}

// ParsingHandler registers to subvisit and parse the scraped HTML.
func (p *PortalNews) ParsingHandler(cd chan<- gachifinder.GachiData) {
	p.c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	p.c.OnResponse(func(r *colly.Response) {
		// fmt.Println(string(r.Body))
	})

	// The headline photo news of left side on news.naver.com
	p.c.OnHTML(".hdline_flick_item", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, e *colly.HTMLElement) {
			link := e.Attr("href")
			if text := e.ChildText("p.hdline_flick_tit"); text != "" {
				// fmt.Println(".hdline_flick_item: Link found:", text, "->", link)

				// Visit link found on page on a new thread(go routine)
				e.Request.Visit(link)
			}
		})
	})

	// The headline news list of right side on news.naver.com
	p.c.OnHTML(".hdline_article_tit", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, e *colly.HTMLElement) {
			link := e.Attr("href")
			if text := strings.TrimSpace(strings.Trim(e.Text, "\n")); text != "" {
				// fmt.Println(".hdline_article_tit: Link found:", text, "->", link)

				// Visit link found on page on a new thread
				e.Request.Visit(link)
			}
		})
	})

	// The entire news except headline on news.naver.com
	p.c.OnHTML(".com_list", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, e *colly.HTMLElement) {
			link := e.Attr("href")
			if text := strings.TrimSpace(strings.Trim(e.Text, "\n")); text != "" {
				// fmt.Println(".com_list: Link found:", text, "->", link)

				// Visit link found on page on a new thread
				e.Request.Visit(link)
			}
		})
	})

	p.c.OnHTML("head", func(e *colly.HTMLElement) {
		if e.Request.URL.Path == "/" {
			return // Skip if called from the root domain like "news.naver.com"
		}

		emitData := gachifinder.GachiData{
			Timestamp: p.timestamp,
			ShortCutIconURL: e.ChildAttr("link[rel='shortcut icon']", "href"),
			Title: e.ChildAttr("meta[name='twitter:title']", "content"),
			URL: e.ChildAttr("meta[property='og:url']", "content"),
			ImageURL: e.ChildAttr("meta[name='twitter:image']", "content"),
			Creator: e.ChildAttr("meta[name='twitter:creator']", "content"),
			Description: e.ChildAttr("meta[name='twitter:description']", "content"),
		}

		cd <- emitData
	})

	p.c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
}