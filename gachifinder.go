package gachifinder

import (
	"time"
)

// ParsingHandler ...
type ParsingHandler func(chan<- GachiData)

// Scraper interface is a crawling actor.
type Scraper interface {
	Do(ParsingHandler, chan<- GachiData, chan<- bool)
	ParsingHandler(chan<- GachiData)
}

// Emitter interface to sent or write the data to the targets.
type Emitter interface {
	// Connect to the Emitter; connect is only called once when the plugin starts.
	Connect() error
	// Close any connections to the Emitter. Close is called once when the output
	// is shutting down. Close will not be called until all writes have finished,
	// and Write() will not be called once Close() has been, so locking is not
	// necessary.
	Close()
	// Write takes in group of points to be written to the Emitter
	Write(<-chan GachiData, <-chan bool) error
}

// GachiData is contents to collect data by scraper.
type GachiData struct {
	Timestamp		time.Time
	Creator			string
	Title			string
	Description		string
	URL				string
	ShortCutIconURL	string
	ImageURL		string
}