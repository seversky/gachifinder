# GachiFinder

This project is an agent for scraping, parsing and writing to some storages.<br />
Firstly, it has been scraped the news on the portal(eg, Naver/Daum) in Korea.

## Workflow(Pipeline)

```
Target page -------------------------|
  |- 1'st Sub -----------------------|
       |- 2'st Sub ------------------|
             |- N'st visitable ------|=> asynchronous scraper(crawler) module
       |- ... -----------------------|                  |
  |- 1'st Sub -----------------------|                  |=> asynchronous emitter(store or relay) module
  |- ... ----------------------------|
```

### **Scraper Modules**

Actually, this is an extendable parsing handler.<br />
Data crawling do via the [Scraper](https://github.com/seversky/gachifinder/blob/master/scrape/scrape.go) interface.

[naver_news_headline](https://github.com/seversky/gachifinder/blob/master/scrape/naver_news_headline.go)<br />
[daum_news_headline](https://github.com/seversky/gachifinder/blob/master/scrape/daum_news_headline.go)

### **Emitter Modules**

This is an extendable output module via [Emitter](https://github.com/seversky/gachifinder/blob/master/gachifinder.go) interface<br />

[elasticsearch](https://github.com/seversky/gachifinder/blob/master/emit/elasticsearch.go)

## How to build

### **Preparing dependencies**

- GachiFinder works with go 1.15+(on current build version).
- In the case of **Windows**, it couldn't be invoked "make" command,<br>
  So you need to download and install [GNUMake](http://gnuwin32.sourceforge.net/packages/make.htm) for windows.
  - [Direct download (click)](http://gnuwin32.sourceforge.net/downlinks/make.php)

### **Run from the source code**

#### Tested Support OS : Linux, MacOSX(darwin), Windows

```bash
# If you're on Windows, run "Git Bash" and type the followings.

$ git clone https://github.com/seversky/gachifinder.git
$ cd gachifinder
$ make all # or one of "windows", "darwin" and "linux".
```

If well done, you can see the binary.

```bash
$ cd $GACHIFINDER_FOLDER/cmd/gachifinder/windows
$ ls
gachifinder.exe
```

You can run it refers to the help options.

```bash
$ ./gachifinder.exe -h
Usage:
  C:\Users\...\go\src\github.com\seversky\gachifinder\cmd\gachifinder\windows\gachifinder.exe [OPTIONS]

Options for GachiFinder

Application Options:
  /c, /config_path:  Path To configure (default: ../config/gachifinder.yml)
                     [%CONFIG_PATH%]
  /t, /test          To test for crawling via a scraper only
                     (Without an emitter module)
  /v, /version       Show GachiFinder version and git revision id

Help Options:
  /?                 Show this help message
  /h, /help          Show this help message
```

### **Options: gachifinder.yml**

```yaml
global:
  max_used_cores: 0 # if zero(0), all cores used.
  interval: 5 # Crawing interval(unit: min)

  log:
    log_level: debug # one of trace, debug, info, warn[ing], error, fatal or panic
    stdout: true
    format: 'text' # one of "text" or "json"
    # go_time_format: '2006-01-02T15:04:05.999Z07:00' # default=RFC3339, refer to https://golang.org/src/time/format.go
    force_colors: true

    log_path: './log/gachifinder.log'
    max_size: 50 # Max megabytes before log is rotated
    max_age: 7 # Max number of days to retain log files
    max_backups: 3 # Max number of old log files to keep
    compress: true

scraper:
  visit_domains:
    - https://news.naver.com
    - https://news.daum.net
  # allowed_domains:
  #   - https://news.naver.com
  #   - https://news.daum.net
  user_agent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36'
  max_depth_to_visit: 1

  async: true # Async turns on asynchronous network communication. Use Collector.Wait() to be sure all requests have been finished.

  parallelism: 20 # The number of the maximum allowed concurrent requests of the matching domains.
  delay: 1 # The duration to wait before creating a new request to the matching domains.(unit: sec)
  random_delay: 5 # The extra randomized duration to wait added to delay before creating a new request.(unit: sec)

  consumer_queue_threads: 2 # The number of consumer queue threads
  consumer_queue_max_size: 10 # Max size of consumer queue

emitter:
  elasticsearch:
    hosts:
      - http://elasticsearch:9200
    username: elastic
    password: changeme
```

### **Simple test for scraping**

```bash
# Go into the folder where the built binary of gachifinder is.
$ cd $GACHIFINDER_FOLDER/cmd/gachifinder/windows
$ ./gachifinder.exe -t
INFO[2021-05-28T15:13:41+09:00] Show All Configurations                       config.Emitter.Elasticsearch.Hosts="[http://elasticsearch:9200]" config.Emitter.Elasticsearch.Password=changem
e config.Emitter.Elasticsearch.Username=elastic config.Global.Interval=5 config.Global.Log.Compress=true config.Global.Log.ForceColors=true config.Global.Log.Format=text
config.Global.Log.GoTimeFormat= config.Global.Log.LogLevel=debug config.Global.Log.LogPath=./log/gachifinder.log config.Global.Log.MaxAge=7 config.Global.Log.MaxBackups=3
 config.Global.Log.MaxSize=50 config.Global.Log.Stdout=true config.Global.MaxUsedCores=0 config.Scraper.AllowedDomains="[]" config.Scraper.Async=true config.Scraper.
ConsumerQueueMaxSize=10 config.Scraper.ConsumerQueueThreads=2 config.Scraper.Delay=1 config.Scraper.MaxDepthToVisit=1 config.Scraper.Parallelism=20 config.Scraper.RandomD
elay=5 config.Scraper.UserAgent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36" config.Scraper.VisitDomains="[https://news.n
aver.com https://news.daum.net]"
INFO[2021-05-28T15:13:41+09:00] Application Initializing                      1-GO Runtime Version=go1.16.3 2-System Arch=amd64 3-GachiFider version=0.1.0 4-GachiFider revisi
on number=c44ad459 5-Number of used CPUs=8
INFO[2021-05-28T14:56:34+09:00] Begin crawling
INFO[2021-05-28T14:56:34+09:00] visiting https://news.daum.net
INFO[2021-05-28T14:56:34+09:00] visiting https://news.naver.com
...
(omit)
```

## Using Docker with Elasticsearch and Kibana for a test environment on Linux.

_I suppose to be already installed Docker and Docker-compose therefore I don't handle installing those here._

To run Elasticsearch and Kibana, just go ahead below.

```bash
$ cd $GACHIFINDER_FOLDER/docker
$ docker-compose up -d --build
```

If Elasticsearch account **_has been changed_**, you need to type into [docker/kibana/kibana.yml](https://github.com/seversky/gachifinder/blob/master/docker/kibana/kibana.yml)
