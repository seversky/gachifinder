package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/seversky/gachifinder"
	"github.com/seversky/gachifinder/scrape"
)

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

	// This is the temporary routine run every 5 minutes.
	// To do: I'll apply one of some scheduler modules. eg, github.com/go-co-op/gocron.
	for {
		collectedData := make([]gachifinder.GachiData, 10, 20)
		p.Do(p.ParsingHandler, collectedData)
		time.Sleep(5 * time.Minute)
	}
}
