package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jessevdk/go-flags"
	"github.com/seversky/gachifinder"
	"github.com/seversky/gachifinder/emit"
	"github.com/seversky/gachifinder/scrape"
)

const esURL = "http://localhost:9200"

// Options : Cli option Flags
type Options struct {
	Daemon bool `short:"d" long:"daemon" description:"To run it daemon mode."`
	Config flags.Filename `short:"c" long:"config" default:"../config/gachifinder.yml" env:"CONFIG" description:"Path To configure."`
	Test bool `short:"t" long:"test" description:"To test for crawling via a scraper only.\n(Without an emitter module)\nNOTE: Cannot Run with '-d'(daemon)"`
}

var options Options

func main() {
	var parser = flags.NewParser(&options, flags.Default)
	parser.ShortDescription = `GachiFinder`
	parser.LongDescription = `Options for GachiFinder`

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		fmt.Println("E! The program has been anomaly exited. Exit code =", code)
		os.Exit(code)
	}

	runtime.GOMAXPROCS(1)
	fmt.Println("Starting gachifinder with", runtime.GOMAXPROCS(0), "core(s).")

	var sc scrape.Scrape
	sc.VisitDomains = []string {
		"https://" + scrape.NaverNews,
		"https://" + scrape.DaumNews,
	}
	// sc.AllowedDomains = []string {
	// 	"news.naver.com",
	// 	"news.daum.net",
	// 	"news.v.daum.net",
	// }

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
