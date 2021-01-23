package emit

import (
	"fmt"
	"time"
	"strconv"
	"context"
	"strings"
	"sync"

	"github.com/olivere/elastic/v7"
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
	URLs                []string

	// Unexported ...
	client *elastic.Client
	majorReleaseNumber  int
}

// Connect to Elasticsearch & Create index.
func (e *Elasticsearch) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	client, err := elastic.NewClient(
		elastic.SetURL(e.URLs...),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10 * time.Second),
		elastic.SetGzip(true),
		)
	if err != nil {
		return err
	}

	// check for ES version on first node.
	esVersion, err := client.ElasticsearchVersion(e.URLs[0])

	if err != nil {
		fmt.Println("Elasticsearch version check failed:", err)
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

	fmt.Println("I! Elasticsearch version: " + esVersion)
	fmt.Println("I! Elasticsearch major version number:", majorReleaseNumber)

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
func (e *Elasticsearch) Write(cd <-chan gachifinder.GachiData) error {
	var wg sync.WaitGroup
	bulkRequest := e.client.Bulk()

	wg.Add(1)
	go func () {
		for data := range cd {
			m := make(map[string]interface{})
			m["@timestamp"] 	= data.Timestamp
			m["creator"] 		= data.Creator
			m["title"] 			= data.Title
			m["description"] 	= data.Description
			m["url"] 			= data.URL
			m["short_icon_url"]	= data.ShortCutIconURL
			m["image_url"] 		= data.ImageURL

			br := elastic.NewBulkIndexRequest().Index(templateName).Doc(m)
			bulkRequest.Add(br)
		}
		wg.Done()
	}()
	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	res, err := bulkRequest.Do(ctx)

	if err != nil {
		return fmt.Errorf("Error sending bulk request to Elasticsearch: %s", err)
	}

	if res.Errors {
		for id, err := range res.Failed() {
			fmt.Printf("E! Elasticsearch indexing failure, id: %d, error: %s, caused by: %s, %s", id, err.Error.Reason, err.Error.CausedBy["reason"], err.Error.CausedBy["type"])
		}
		return fmt.Errorf("W! Elasticsearch failed to index %d metrics", len(res.Failed()))
	}

	return nil
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