package gachifinder

// Config is the options of gachifinder.
type Config struct {
	Global struct {
		MaxUsedCores 	int 	`yaml:"max_used_cores"`
		Interval		int		`yaml:"interval"`
		Log struct {
			LogLevel		string 	`yaml:"log_level"`
			Stdout			bool 	`yaml:"stdout"`
			Format			string 	`yaml:"format"`
			ForceColors		bool	`yaml:"force_colors"`
			GoTimeFormat	string 	`yaml:"go_time_format"`
			LogPath			string	`yaml:"log_path"`
			MaxSize			int 	`yaml:"max_size"`
			MaxAge			int 	`yaml:"max_age"`
			MaxBackups		int 	`yaml:"max_backups"`
			Compress		bool 	`yaml:"compress"`
		}`yaml:"log"`
	} `yaml:"global"`

	Scraper struct {
		VisitDomains			[]string 	`yaml:"visit_domains"`
		AllowedDomains			[]string 	`yaml:"allowed_domains"`
		UserAgent				string 		`yaml:"user_agent"`
		MaxDepthToVisit			int			`yaml:"max_depth_to_visit"`
		Async					bool		`yaml:"async"`
		Parallelism				int			`yaml:"parallelism"`
		Delay					int			`yaml:"delay"`
		RandomDelay				int			`yaml:"random_delay"`
		ConsumerQueueThreads	int			`yaml:"consumer_queue_threads"`
		ConsumerQueueMaxSize	int			`yaml:"consumer_queue_max_size"`
	} `yaml:"scraper"`

	Emitter struct {
		Elasticsearch struct {
			Hosts		[]string	`yaml:"hosts"`
			Username	string		`yaml:"username"`
			Password	string		`yaml:"password"`
		}
	}
}

// Types of Config.Global.Log.Format
const (
	JSON = "json"
	TEXT = "text"
) 

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