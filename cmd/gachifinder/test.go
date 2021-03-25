package main

import (
	"log"
	"sync"

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
			log.Println("I! The number of the collected data:", length)
			for _, data := range emitData {
				log.Println("")
				log.Println("I! Timestamp =", data.Timestamp)
				log.Println("I! VisitHost =", data.VisitHost)
				log.Println("I! Creator =", data.Creator)
				log.Println("I! Title =", data.Title)
				log.Println("I! Description =", data.Description)
				log.Println("I! URL =", data.URL)
				log.Println("I! ShortCutIconURL =", data.ShortCutIconURL)
				log.Println("I! ImageURL =", data.ImageURL)
			}
		} else {
			log.Println("W! There is not any collected data")
		}

		wg.Done()
	}()
	wg.Wait()
}