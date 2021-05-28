package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jessevdk/go-flags"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/seversky/gachifinder"
)

// Options : Cli option Flags
type Options struct {
	ConfigPath flags.Filename `short:"c" long:"config_path" default:"../config/gachifinder.yml" env:"CONFIG_PATH" description:"Path To configure"`
	ScrapeTest bool `short:"t" long:"test" description:"To test for crawling via a scraper only\n(Without an emitter module)"`
	ShowVersion bool   `short:"v" long:"version" description:"Show GachiFinder version and git revision id"`
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
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				os.Exit(0)
			}
		}
		return config, fmt.Errorf("The program has been anomaly exited. Exit code = %d", 1)
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

	err = setLogger(&config)
	if err != nil {
		return config, err
	}

	logger.WithFields(logger.Fields{
		"config.Global.MaxUsedCores": config.Global.MaxUsedCores,
		"config.Global.Interval": config.Global.Interval,
		"config.Global.Log.LogLevel": config.Global.Log.LogLevel,
		"config.Global.Log.Stdout": config.Global.Log.Stdout,
		"config.Global.Log.Format": config.Global.Log.Format,
		"config.Global.Log.ForceColors": config.Global.Log.ForceColors,
		"config.Global.Log.GoTimeFormat": config.Global.Log.GoTimeFormat,
		"config.Global.Log.LogPath": config.Global.Log.LogPath,
		"config.Global.Log.MaxSize": config.Global.Log.MaxSize,
		"config.Global.Log.MaxAge": config.Global.Log.MaxAge,
		"config.Global.Log.MaxBackups": config.Global.Log.MaxBackups,
		"config.Global.Log.Compress": config.Global.Log.Compress,
		"config.Scraper.VisitDomains": config.Scraper.VisitDomains,
		"config.Scraper.AllowedDomains": config.Scraper.AllowedDomains,
		"config.Scraper.UserAgent": config.Scraper.UserAgent,
		"config.Scraper.MaxDepthToVisit": config.Scraper.MaxDepthToVisit,
		"config.Scraper.Async": config.Scraper.Async,
		"config.Scraper.Parallelism": config.Scraper.Parallelism,
		"config.Scraper.Delay": config.Scraper.Delay,
		"config.Scraper.RandomDelay": config.Scraper.RandomDelay,
		"config.Scraper.ConsumerQueueThreads": config.Scraper.ConsumerQueueThreads,
		"config.Scraper.ConsumerQueueMaxSize": config.Scraper.ConsumerQueueMaxSize,
		"config.Emitter.Elasticsearch.Hosts": config.Emitter.Elasticsearch.Hosts,
		"config.Emitter.Elasticsearch.Username": config.Emitter.Elasticsearch.Username,
		"config.Emitter.Elasticsearch.Password": config.Emitter.Elasticsearch.Password,
	}).Info("Show All Configurations")

	return config, nil
}