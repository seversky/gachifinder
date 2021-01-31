package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/seversky/gachifinder"
	"github.com/seversky/gachifinder/scrape"
	"github.com/seversky/gachifinder/emit"
)

const esURL = "http://localhost:9200"

func main() {
	runtime.GOMAXPROCS(1)
	fmt.Println("Starting gachifinder with", runtime.GOMAXPROCS(0), "core(s).")

	var p scrape.PortalNews
	p.VisitDomains = []string {
		"https://news.naver.com/",
		// "https://news.daum.net/",
	}
	p.AllowedDomains = []string {
		"news.naver.com",
		"news.naver.com/main",
		"news.daum.net",
		"news.v.daum.net/v",
	}

	var s gachifinder.Scraper = &p

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
		cd := s.Do(s.ParsingHandler)

		err = em.Write(cd)
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
