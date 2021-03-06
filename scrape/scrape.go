package scrape

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
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
	VisitDomains	[]string
	AllowedDomains	[]string

	// Unexport ...
	c 			*colly.Collector	// Will be assigned by inside Do func.
	timestamp 	string
}

// Do creates colly.collector and queue, and then do and wait till done
func (s *Scrape) Do(fs []ParsingHandler) (<-chan gachifinder.GachiData) {
	// Record the beginning time.
	s.timestamp = time.Now().UTC().Format("2006-01-02T15:04:05")
	fmt.Println("I! It gets begun at", time.Now())

	dc := make(chan gachifinder.GachiData)

	go func () {
		// Instantiate default collector
		s.c = colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"),
			colly.Async(true),
			colly.MaxDepth(1),
			colly.AllowedDomains(s.AllowedDomains...),
		)

		s.c.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: 20,
			Delay: time.Second,
			RandomDelay: 5 * time.Second,
		})

		// create a request queue with 1 consumer threads
		q, err := queue.New(
			1, // Number of consumer threads
			&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
		)
		if err != nil {
			fmt.Println("Creating Queue is Failed:", err)
			panic(err)
		}

		for _, url := range s.VisitDomains {
			err := q.AddURL(url)
			if err != nil {
				fmt.Println("Adding url into the queue is Failed:", err)
				panic(err)
			}
		}

		// Common handlers
		s.c.OnRequest(func(r *colly.Request) {
			fmt.Println("visiting", r.URL)
		})
	
		s.c.OnResponse(func(r *colly.Response) {
			// fmt.Println(string(r.Body))
		})

		s.c.OnError(func(r *colly.Response, err error) {
			fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		})

		// Specified Parse handlers.
		for _, f := range fs {
			f(dc, s)
		}

		// Consume URLs.
		err = q.Run(s.c)
		if err != nil {
			fmt.Println("Running the queue is Failed:", err)
			panic(err)
		}
		// Wait for the crawling to complete.
		s.c.Wait()

		close(dc)
	}()

	return dc
}