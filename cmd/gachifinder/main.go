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
		done := make(chan bool)
		cd := make(chan gachifinder.GachiData)

		go p.Do(p.ParsingHandler, cd, done)

		emitData := make([]gachifinder.GachiData, 0, 20)
		for c := true; c;{
			select {
			case data := <-cd:
				emitData = append(emitData, data)
			case <-done:
				c = false
			}
		}

		length := len(emitData)
		if length > 0 {
			fmt.Println(length)
			for _, data := range emitData {
				fmt.Println(data.Timestamp)
				fmt.Println(data.Creator)
				fmt.Println(data.Title)
				fmt.Println(data.Description)
				fmt.Println(data.URL)
				fmt.Println(data.ShortCutIconURL)
				fmt.Println(data.ImageURL)
			}
		} else {
			fmt.Println("There is not any collected data")
		}

		time.Sleep(5 * time.Minute)
	}
}
