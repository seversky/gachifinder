package scrape

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	logger "github.com/sirupsen/logrus"

	"github.com/seversky/gachifinder"
)

// ParsingHandler ...
type ParsingHandler func(chan<- gachifinder.GachiData, *Scrape)

// Scraper interface is a crawling actor.
type Scraper interface {
	// Do is a producer in a part of a pipeline
	Do([]ParsingHandler) (<-chan gachifinder.GachiData)
}

var _ Scraper = &Scrape{}

// Scrape struct.
type Scrape struct {
	Config *gachifinder.Config

	// Unexport ...
	c 			*colly.Collector	// Will be assigned by inside Do func.
	timestamp 	string
}

// Do creates colly.collector and queue, and then do and wait till done
func (s *Scrape) Do(fs []ParsingHandler) (<-chan gachifinder.GachiData) {
	// Record the beginning time.
	s.timestamp = time.Now().UTC().Format("2006-01-02T15:04:05")
	logger.Println("I! It gets begun at", time.Now())

	dc := make(chan gachifinder.GachiData)

	go func () {
		// Instantiate default collector
		s.c = colly.NewCollector(
			colly.UserAgent(s.Config.Scraper.UserAgent),
			colly.MaxDepth(s.Config.Scraper.MaxDepthToVisit),
			colly.AllowedDomains(s.Config.Scraper.AllowedDomains...),
		)

		s.c.Async = s.Config.Scraper.Async

		s.c.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: s.Config.Scraper.Parallelism,
			Delay: time.Duration(s.Config.Scraper.Delay) * time.Second ,
			RandomDelay: time.Duration(s.Config.Scraper.RandomDelay) * time.Second,
		})

		// create a request queue with 1 consumer threads
		q, err := queue.New(
			s.Config.Scraper.ConsumerQueueThreads, // Number of consumer threads
			&queue.InMemoryQueueStorage{MaxSize: s.Config.Scraper.ConsumerQueueMaxSize}, // Use default queue storage
		)
		if err != nil {
			logger.Println("E! Creating Queue is Failed:", err)
			logger.Fatalln("E! error:", err)
		}

		for _, url := range s.Config.Scraper.VisitDomains {
			err := q.AddURL(url)
			if err != nil {
				logger.Println("E! Adding url into the queue is Failed:", err)
				logger.Fatalln("E! error:", err)
			}
		}

		// Common handlers
		s.c.OnRequest(func(r *colly.Request) {
			logger.Println("I! visiting", r.URL)
		})
	
		s.c.OnResponse(func(r *colly.Response) {
			// logger.Println(string(r.Body))
		})

		s.c.OnError(func(r *colly.Response, err error) {
			logger.Println("E! Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		})

		// Specified Parse handlers.
		for _, f := range fs {
			f(dc, s)
		}

		// Consume URLs.
		err = q.Run(s.c)
		if err != nil {
			logger.Println("E! Running the queue is Failed:", err)
			logger.Fatalln("E! error:", err)
		}
		// Wait for the crawling to complete.
		s.c.Wait()

		close(dc)
	}()

	return dc
}