package scrape

import (
	"github.com/gocolly/colly/v2"
	logger "github.com/sirupsen/logrus"

	"github.com/seversky/gachifinder"
)

// DaumNews is the root domain for visiting.
const DaumNews = "news.daum.net"

// OnHTMLDaumHeadlineNews registers to subvisit and parse the scraped "new.daum.com" HTML.
func OnHTMLDaumHeadlineNews(dc chan<- gachifinder.GachiData, s *Scrape) {
	s.c.OnHTML("ul.list_issue", func(e *colly.HTMLElement) {
		if e.Request.URL.Host != DaumNews {
			return
		}

		e.ForEach("a[href]", func(_ int, e *colly.HTMLElement) {
			if e.Attr("class") == "link_txt" {
				link := e.Attr("href")
				e.Request.Visit(link)
			}
		})
	})

	s.c.OnHTML("div.box_headline", func(e *colly.HTMLElement) {
		if e.Request.URL.Host != DaumNews {
			return
		}

		e.ForEach("a[href]", func(_ int, e *colly.HTMLElement) {
			link := e.Attr("href")
			e.Request.Visit(link)
		})
	})

	s.c.OnHTML("head", func(e *colly.HTMLElement) {
		if e.Request.URL.Host != "news.v.daum.net" || len(e.Request.URL.Path) < 2 {
			return // Skip if called from the root domain like "news.daum.com"
		}

		press, isExist := e.DOM.Siblings().Find("div.head_view").Find("img.thumb_g").Attr("alt")
		if isExist {
			url := e.Request.URL.Scheme + "://" + e.Request.URL.Host + e.Request.URL.Path
			emitData := gachifinder.GachiData{
				Timestamp: s.timestamp,
				VisitHost: DaumNews,
				ShortCutIconURL: "https:" + e.ChildAttr("link[rel='shortcut icon']", "href"),
				Title: e.ChildAttr("meta[property='og:title']", "content"),
				URL: url,
				ImageURL: e.ChildAttr("meta[property='og:image']", "content"),
				Creator: press,
				Description: e.ChildAttr("meta[property='og:description']", "content"),
			}

			dc <- emitData
		} else {
			logger.Println("W! There is not any press name!")
		}
	})
}