# gachifinder

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

## Scraper Modules

Actually, this is an extendable parsing handler.<br />
Data crawling do via the [Scraper](https://github.com/seversky/gachifinder/blob/master/scrape/scrape.go) interface.

[naver_news_headline](https://github.com/seversky/gachifinder/blob/master/scrape/naver_news_headline.go)<br />
[daum_news_headline](https://github.com/seversky/gachifinder/blob/master/scrape/daum_news_headline.go)

## Emitter Modules

This is an extendable output module via [Emitter](https://github.com/seversky/gachifinder/blob/master/gachifinder.go) interface<br />

[elasticsearch](https://github.com/seversky/gachifinder/blob/master/emit/elasticsearch.go)

## However, this is still a prototype!

### To-do list before a release version.

- Add configuration file to optionize.
  - For something like crawling target page, emitter ip, user account, etc.
- Application daemonize.
- Definition for a log format.

## Using Docker for a test environment.

```
git clone https://github.com/seversky/gachifinder.git
cd docker
docker-compose up -d --build
```

If Elasticsearch account **_has been changed_**, you need to type into [docker/kibana/kibana.yml](https://github.com/seversky/gachifinder/blob/master/docker/kibana/kibana.yml)
