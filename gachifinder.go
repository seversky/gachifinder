package gachifinder

// Emitter interface to sent or write the data to the targets.
type Emitter interface {
	// Connect to the Emitter; connect is only called once when the plugin starts.
	Connect() error
	// Close any connections to the Emitter. Close is called once when the output
	// is shutting down. Close will not be called until all writes have finished,
	// and Write() will not be called once Close() has been, so locking is not
	// necessary.
	Close()
	// Write takes in group of points to be written to the Emitter.
	// this is a consumer in a part of a pipeline.
	Write(<-chan GachiData) error
}

// GachiData is contents to collect data by scraper.
type GachiData struct {
	Timestamp		string
	VisitHost		string
	Creator			string
	Title			string
	Description		string
	URL				string
	ShortCutIconURL	string
	ImageURL		string
}