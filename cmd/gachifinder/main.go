package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/seversky/gachifinder"
	"github.com/seversky/gachifinder/scrape"
	"github.com/seversky/gachifinder/emit"
)

const esURL = "http://192.168.56.105:9200"

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
		fmt.Println(err)
	}
	defer em.Close()

	// This is the temporary routine run around every 5 minutes.
	// To do: I'll apply one of some scheduler modules. eg, github.com/go-co-op/gocron.
	for {
		cd := s.Do(s.ParsingHandler)

		err = em.Write(cd)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("I! Crawling success", time.Now())
		time.Sleep(4 * time.Minute)
	}
}
