package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/seversky/gachifinder"
	"gopkg.in/yaml.v2"
)

// Options : Cli option Flags
type Options struct {
	Daemon bool `short:"d" long:"daemon" description:"To run it daemon mode."`
	ConfigPath flags.Filename `short:"c" long:"config_path" default:"../config/gachifinder.yml" env:"CONFIG_PATH" description:"Path To configure."`
	ScrapeTest bool `short:"t" long:"test" description:"To test for crawling via a scraper only.\n(Without an emitter module)\nNOTE: Cannot Run with '-d'(daemon)"`
	ShowVersion bool   `short:"v" long:"version" description:"Show GachiFinder version info"`
}

var options Options

// Build flags
var (
	version = ""
	commit  = ""
)

func setOptions() (gachifinder.Config, error) {
	config := gachifinder.Config{}

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
		return config, fmt.Errorf("E! The program has been anomaly exited. Exit code = %d", code)
	}

	if options.ShowVersion {
		fmt.Printf("gachifinder %s (git commit: %s)\n", version, commit)
		os.Exit(0)
	}

	bytes, err := ioutil.ReadFile(string(options.ConfigPath))
    if err != nil {
		return config, err
    }

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}

	log.Println("I! config.Global.MaxUsedCores =", config.Global.MaxUsedCores)
	log.Println("I! config.Global.Interval =", config.Global.Interval)
	log.Println("I! config.Scraper.VisitDomains =", config.Scraper.VisitDomains)
	log.Println("I! config.Scraper.AllowedDomains =", config.Scraper.AllowedDomains)
	log.Println("I! config.Scraper.UserAgent =", config.Scraper.UserAgent)
	log.Println("I! config.Scraper.MaxDepthToVisit =", config.Scraper.MaxDepthToVisit)
	log.Println("I! config.Scraper.Async =", config.Scraper.Async)
	log.Println("I! config.Scraper.Parallelism =", config.Scraper.Parallelism)
	log.Println("I! config.Scraper.Delay =", config.Scraper.Delay)
	log.Println("I! config.Scraper.RandomDelay =", config.Scraper.RandomDelay)
	log.Println("I! config.Scraper.ConsumerQueueThreads =", config.Scraper.ConsumerQueueThreads)
	log.Println("I! config.Scraper.ConsumerQueueMaxSize =", config.Scraper.ConsumerQueueMaxSize)
	log.Println("I! config.Emitter.Elasticsearch.Hosts =", config.Emitter.Elasticsearch.Hosts)
	log.Println("I! config.Emitter.Elasticsearch.Username =", config.Emitter.Elasticsearch.Username)
	log.Println("I! config.Emitter.Elasticsearch.Password =", config.Emitter.Elasticsearch.Password)

	return config, nil
}