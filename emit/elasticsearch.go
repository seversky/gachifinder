package emit

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
	logger "github.com/sirupsen/logrus"

	"github.com/seversky/gachifinder"
)

const templateName = "gachifinder"
const indexTemplate = `
{
	"index_patterns": [
		"gachifinder*"
	],
	"settings": {
		"number_of_shards": 1,
		"number_of_replicas": 1,
		"index": {
			"analysis": {
				"analyzer": {
					"gachi_analyzer": {
						"type": "custom",
						"tokenizer": "gachi_user_dict",
						"filter": [
							"gachi_posfilter"
						]
					}
				},
				"tokenizer": {
					"gachi_user_dict": {
						"type": "nori_tokenizer",
						"decompound_mode": "mixed",
						"user_dictionary": "userdict_ko.txt"
					}
				},
				"filter": {
					"gachi_posfilter": {
						"type": "nori_part_of_speech",
						"stoptags": [
							"E",
							"EF",
							"EC",
							"IC",
							"J",
							"MAG",
							"MAJ",
							"MM",
							"NA",
							"SP",
							"SSC",
							"SSO",
							"SC",
							"SE",
							"UNA",
							"VSV",
							"VA",
							"VV",
							"VX",
							"XPN",
							"XSA",
							"XSN",
							"XSV"
						]
					}
				},
				"normalizer": {
					"gachi_normalizer": {
						"type": "custom",
						"filter": [
							"lowercase",
							"asciifolding"
						]
					}
				}
			}
		}
	},
	"mappings": {
		"properties": {
			"@timestamp": {
				"type": "date"
			},
			"visit_host": {
				"type": "keyword",
				"normalizer": "gachi_normalizer"
			},
			"creator": {
				"type": "text",
				"analyzer": "gachi_analyzer",
				"fields": {
					"keyword": {
						"type": "keyword",
						"normalizer": "gachi_normalizer"
					}
				}
			},
			"title": {
				"type": "text",
				"analyzer": "gachi_analyzer",
				"fields": {
					"keyword": {
						"type": "keyword",
						"normalizer": "gachi_normalizer"
					}
				}
			},
			"description": {
				"type": "text",
				"analyzer": "gachi_analyzer",
				"fields": {
					"keyword": {
						"type": "keyword",
						"normalizer": "gachi_normalizer"
					}
				}
			},
			"url": {
				"type": "keyword",
				"normalizer": "gachi_normalizer"
			},
			"short_icon_url": {
				"type": "text",
				"index": false
			},
			"image_url": {
				"type": "text",
				"index": false
			}
		}
	}
}`

var _ gachifinder.Emitter = &Elasticsearch{}

// Elasticsearch struct
type Elasticsearch struct {
	Config *gachifinder.Config

	// Unexport ...
	client 				*elastic.Client
	majorReleaseNumber	int
}

// Connect to Elasticsearch & Create index.
func (e *Elasticsearch) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	client, err := elastic.NewClient(
		elastic.SetURL(e.Config.Emitter.Elasticsearch.Hosts...),
		elastic.SetBasicAuth(e.Config.Emitter.Elasticsearch.Username, e.Config.Emitter.Elasticsearch.Password),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10 * time.Second),
		elastic.SetGzip(true),
		)
	if err != nil {
		return err
	}

	// check for ES version on first node.
	esVersion, err := client.ElasticsearchVersion(e.Config.Emitter.Elasticsearch.Hosts[0])

	if err != nil {
		logger.Error("Elasticsearch version check failed")
		return err
	}

	// quit if ES version is not supported.
	majorReleaseNumber, err := strconv.Atoi(strings.Split(esVersion, ".")[0])
	if err != nil {
		return err
	}
	if majorReleaseNumber < 7 {
		return fmt.Errorf("Elasticsearch version not supported: %s", esVersion)
	}

	logger.WithField("Elasticsearch version", esVersion).
		WithField("Elasticsearch major version number", majorReleaseNumber).
		Info("Elasticsearch Info")

	e.client = client
	e.majorReleaseNumber = majorReleaseNumber

	err = e.manageTemplate(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Close to release Elasticsearch Client.
func (e *Elasticsearch) Close() {
	e.client = nil
}

// Write the data into the Elasticsearch.
func (e *Elasticsearch) Write(dc <-chan gachifinder.GachiData) error {
	var wg sync.WaitGroup
	bulkRequest := e.client.Bulk()

	wg.Add(1)
	go func () {
		for data := range dc {
			logger.WithField("1-Timestamp", data.Timestamp).
				WithField("2-VisitHost", data.VisitHost).
				WithField("3-Creator", data.Creator).
				WithField("4-Title", data.Title).
				WithField("5-Description", data.Description).
				WithField("6-URL", data.URL).
				WithField("7-ShortCutIconURL", data.ShortCutIconURL).
				WithField("8-ImageURL", data.ImageURL).
				Debug("Collected data")

			m := make(map[string]interface{})
			m["@timestamp"] 	= data.Timestamp
			m["visit_host"]		= data.VisitHost
			m["creator"] 		= data.Creator
			m["title"] 			= data.Title
			m["description"] 	= data.Description
			m["url"] 			= data.URL
			m["short_icon_url"]	= data.ShortCutIconURL
			m["image_url"] 		= data.ImageURL

			indexName := e.generateIndexName(data.Timestamp)
			br := elastic.NewBulkIndexRequest().Index(indexName).Doc(m)
			bulkRequest.Add(br)
		}
		wg.Done()
	}()
	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	for retry := 1; retry < 4; retry++ {
		res, err := bulkRequest.Do(ctx)

		if err != nil {
			logger.WithField("Retry count", retry).
				WithField("error", err).
				Error("Elasticsearch bulk writing error")
			continue
		} else {
			if res.Errors {
				for id, err := range res.Failed() {
					logger.WithField("id", id).
						WithField("error", err.Error.Reason).
						WithField("reason caused by", err.Error.CausedBy["reason"]).
						WithField("type caused by", err.Error.CausedBy["type"]).
						Error("Elasticsearch indexing failure")
				}
				return fmt.Errorf("Elasticsearch failed to index %d metrics", len(res.Failed()))
			}
			return nil
		}
	}

	return fmt.Errorf("Retry counts are exceeded")
}

func (e *Elasticsearch) manageTemplate(ctx context.Context) error {
	templateExists, errExists := e.client.IndexTemplateExists(templateName).Do(ctx)
	if errExists != nil {
		return fmt.Errorf("Elasticsearch template check failed, template name: '%s', error: %s", templateName, errExists)
	}

	if !templateExists {
		_, errCreateTemplate := e.client.IndexPutTemplate(templateName).BodyString(indexTemplate).Do(ctx)

		if errCreateTemplate != nil {
			return fmt.Errorf("Elasticsearch failed to create index template '%s' : %s", templateName, errCreateTemplate)
		}
		templateExists, errExists := e.client.IndexTemplateExists(templateName).Do(ctx)
		if errExists != nil {
			return fmt.Errorf("Elasticsearch template check failed, template name: '%s', error: %s", templateName, errExists)
		}
		if !templateExists {
			return fmt.Errorf("Failed to create the template '%s'", templateName)
		}
	}

	return nil
}

func (e *Elasticsearch) generateIndexName(timestamp string) string {
	n := timestamp[:10]
	return fmt.Sprintf("%s_%s", templateName, n)
}