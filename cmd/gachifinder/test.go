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
			logger.Println("I! The number of the collected data:", length)
			for _, data := range emitData {
				logger.Println("")
				logger.Println("I! Timestamp =", data.Timestamp)
				logger.Println("I! VisitHost =", data.VisitHost)
				logger.Println("I! Creator =", data.Creator)
				logger.Println("I! Title =", data.Title)
				logger.Println("I! Description =", data.Description)
				logger.Println("I! URL =", data.URL)
				logger.Println("I! ShortCutIconURL =", data.ShortCutIconURL)
				logger.Println("I! ImageURL =", data.ImageURL)
			}
		} else {
			logger.Println("W! There is not any collected data")
		}

		wg.Done()
	}()
	wg.Wait()
}