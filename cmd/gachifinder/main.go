package main

import (
	"os"
	"runtime"
	"time"

	"github.com/go-co-op/gocron"
	logger "github.com/sirupsen/logrus"

	"github.com/seversky/gachifinder"
	"github.com/seversky/gachifinder/emit"
	"github.com/seversky/gachifinder/scrape"
)

func main() {
	// Set command options and config options
	config, err := setOptions()
	if err != nil {
		logger.Fatalln("E! error:", err)
	}

	// Set used core(s)
	if config.Global.MaxUsedCores == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(config.Global.MaxUsedCores)
	}
	logger.WithFields(logger.Fields{
		"1:GO Runtime Version": runtime.Version(),
		"2:System Arch": runtime.GOARCH,
		"3:GachiFider version": version,
		"4:GachiFider revision number": commit,
		"5:Number of used CPUs": runtime.GOMAXPROCS(0),
	}).Info("Application Initializing")
	
	os.Exit(0) // Jack: test

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
		logger.Fatalln("E! error:", err)
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
			logger.Println(err)
			logger.Println("E! Crawling is failed at", time.Now())
		} else {
			logger.Println("I! Crawling is done successfully at", time.Now())
		}
		_, tNext := js.NextRun()
		logger.Println("I! It'll get begun at", tNext)
	})
	if errJs != nil {
        logger.Fatalln("E! error:", err)
    }

	js.StartBlocking()
}
