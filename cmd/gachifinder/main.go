package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/seversky/gachifinder"
	"github.com/seversky/gachifinder/emit"
	"github.com/seversky/gachifinder/scrape"
)

const esURL = "http://localhost:9200"

func main() {
	runtime.GOMAXPROCS(1)
	fmt.Println("Starting gachifinder with", runtime.GOMAXPROCS(0), "core(s).")

	var sc scrape.Scrape
	sc.VisitDomains = []string {
		"https://" + scrape.NaverNews,
		"https://" + scrape.DaumNews,
	}
	sc.AllowedDomains = []string {
		"news.naver.com",
		"news.daum.net",
		"news.v.daum.net",
	}

	var s scrape.Scraper = &sc

	e := &emit.Elasticsearch {
		URLs: []string{esURL},
	}

	var em gachifinder.Emitter = e

	err := em.Connect()
	if err != nil {
		panic(err)
	}
	defer em.Close()

	// defines a new scheduler that schedules and runs jobs
	js := gocron.NewScheduler(time.Local)
	_, errJs := js.Every(5).Minutes().Do(func() {
		fs := []scrape.ParsingHandler {
			scrape.OnHTMLNaverHeadlineNews,
			scrape.OnHTMLDaumHeadlineNews,
		}
		dc := s.Do(fs)

		err = em.Write(dc)
		if err != nil {
			fmt.Println(err)
			fmt.Println("E! Crawling is failed at", time.Now())
		} else {
			fmt.Println("I! Crawling is done successfully at", time.Now())
		}
		_, tNext := js.NextRun()
			fmt.Println("I! It'll get begun at", tNext)
	})
	if errJs != nil {
        panic(err)
    }

	js.StartBlocking()
}
