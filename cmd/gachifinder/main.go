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
		logger.WithField("error", err).Fatalln("Option or configuration fail")
	}

	// Set used core(s)
	if config.Global.MaxUsedCores == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(config.Global.MaxUsedCores)
	}
	logger.WithField("1-GO Runtime Version", runtime.Version()).
		WithField("2-System Arch", runtime.GOARCH).
		WithField("3-GachiFider version", version).
		WithField("4-GachiFider revision number", commit).
		WithField("5-Number of used CPUs", runtime.GOMAXPROCS(0)).
		Info("Application Initializing")
	
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
		logger.WithField("error", err).Fatalln("Emitter fail")
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
			logger.WithField("error", err).Error("Crawling is failed")
		} else {
			logger.Info("Crawling is done successfully")
		}
		_, tNext := js.NextRun()
		logger.WithField("Next running time", tNext).Info("Crawling time")
	})
	if errJs != nil {
        logger.WithField("error", err).Fatalln("Job scheduling fail")
    }

	js.StartBlocking()
}
