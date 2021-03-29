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

## However, this is still a prototype!

### To-do list before a release version.

- Application daemonize.
- Definition for a log format.

## How to build

### **Preparing dependencies**

- GachiFinder works with go 1.15+(on current build version).
- In the case of **Windows**, it couldn't be invoked "make" command,<br>
  So you need to download and install [GNUMake](http://gnuwin32.sourceforge.net/packages/make.htm) for windows.
  - [Direct download (click)](http://gnuwin32.sourceforge.net/downlinks/make.php)

### **Run from the source code**

```bash
# If you're on Windows, run "Git Bash" and type the followings.

$ go get github.com/seversky/gachifinder
$ cd $GOPATH/src/github.com/seversky/gachifinder
$ make all # or one of "windows", "darwin" and "linux".
```

If well done, you can see the binary.

```bash
$ cd cmd/gachifinder/windows
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
  /d, /daemon        To run it daemon mode, but not supported yet!
  /c, /config_path:  Path To configure (default: ../config/gachifinder.yml)
                     [%CONFIG_PATH%]
  /t, /test          To test for crawling via a scraper only
                     (Without an emitter module)
                     NOTE: Cannot Run with '-d'(daemon)
  /v, /version       Show GachiFinder version and git revision id

Help Options:
  /?                 Show this help message
  /h, /help          Show this help message
```

## Using Docker with Elasticsearch and Kibana for a test environment on Linux.

_I suppose to be already installed Docker and Docker-compose therefore I don't handle installing those here._

To run Elasticsearch and Kibana, just go ahead below.

```bash
$ cd $GOPATH/src/github.com/seversky/gachifinder/docker
$ docker-compose up -d --build
```

If Elasticsearch account **_has been changed_**, you need to type into [docker/kibana/kibana.yml](https://github.com/seversky/gachifinder/blob/master/docker/kibana/kibana.yml)
