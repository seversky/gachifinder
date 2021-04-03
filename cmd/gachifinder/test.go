package main

import (
	"sync"

	logger "github.com/sirupsen/logrus"

	"github.com/seversky/gachifinder"
	"github.com/seversky/gachifinder/scrape"
)

func scrapeTest(config *gachifinder.Config) {
	var sc scrape.Scrape = scrape.Scrape {
		Config: config,
	}

	var s scrape.Scraper = &sc

	fs := []scrape.ParsingHandler {
		scrape.OnHTMLNaverHeadlineNews,
		scrape.OnHTMLDaumHeadlineNews,
	}
	dc := s.Do(fs)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		emitData := make([]gachifinder.GachiData, 0, 20)
		for data := range dc {
			emitData = append(emitData, data)
		}

		length := len(emitData)
		if length > 0 {
			logger.WithField("The number of the collected data", length).Info("Crawling finished")
			for _, data := range emitData {
				logger.WithField("1-Timestamp", data.Timestamp).
					WithField("2-VisitHost", data.VisitHost).
					WithField("3-Creator", data.Creator).
					WithField("4-Title", data.Title).
					WithField("5-Description", data.Description).
					WithField("6-URL", data.URL).
					WithField("7-ShortCutIconURL", data.ShortCutIconURL).
					WithField("8-ImageURL", data.ImageURL).
					Info("Collected data")
			}
		} else {
			logger.Warn("There is not any collected data")
		}

		wg.Done()
	}()
	wg.Wait()
}