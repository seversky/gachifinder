package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/seversky/gachifinder"
	"github.com/seversky/gachifinder/emit"
	"github.com/seversky/gachifinder/scrape"
)

func main() {
	// Set command options and config options
	config, err := setOptions()
	if err != nil {
		log.Fatalln("E! error:", err)
	}

	// Set used core(s)
	if config.Global.MaxUsedCores == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(config.Global.MaxUsedCores)
	}
	log.Println("I! Starting gachifinder with", runtime.GOMAXPROCS(0), "core(s).")

	if options.ScrapeTest {
		scrapeTest(&config)
		os.Exit(0)
	}

	var sc scrape.Scrape = scrape.Scrape {
		Config: &config,
	}

	var s scrape.Scraper = &sc

	e := &emit.Elasticsearch {
		Config: &config,
	}

	var em gachifinder.Emitter = e

	err = em.Connect()
	if err != nil {
		log.Fatalln("E! error:", err)
	}
	defer em.Close()

	// defines a new scheduler that schedules and runs jobs
	js := gocron.NewScheduler(time.Local)
	_, errJs := js.Every(uint64(config.Global.Interval)).Minutes().Do(func() {
		fs := []scrape.ParsingHandler {
			scrape.OnHTMLNaverHeadlineNews,
			scrape.OnHTMLDaumHeadlineNews,
		}
		dc := s.Do(fs)

		err = em.Write(dc)
		if err != nil {
			log.Println(err)
			log.Println("E! Crawling is failed at", time.Now())
		} else {
			log.Println("I! Crawling is done successfully at", time.Now())
		}
		_, tNext := js.NextRun()
		log.Println("I! It'll get begun at", tNext)
	})
	if errJs != nil {
        log.Fatalln("E! error:", err)
    }

	js.StartBlocking()
}
