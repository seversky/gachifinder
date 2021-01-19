package emit

import (
	"fmt"
	"time"
	"strconv"
	"context"
	"strings"

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
					"nori_analyzer": {
						"type": "custom",
						"tokenizer": "nori_user_dict",
						"filter": [
							"my_posfilter"
						]
					}
				},
				"tokenizer": {
					"nori_user_dict": {
						"type": "nori_tokenizer",
						"decompound_mode": "mixed",
						"user_dictionary": "userdict_ko.txt"
					}
				},
				"filter": {
					"my_posfilter": {
						"type": "nori_part_of_speech",
						"stoptags": [
							"E",
							"IC",
							"J",
							"MAG",
							"MAJ",
							"MM",
							"SP",
							"SSC",
							"SSO",
							"SC",
							"SE",
							"XPN",
							"XSA",
							"XSN",
							"XSV",
							"UNA",
							"NA",
							"VSV"
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
				"type": "text"
			},
			"title": {
				"type": "text",
				"analyzer": "nori_analyzer"
			},
			"description": {
				"type": "text",
				"analyzer": "nori_analyzer"
			},
			"url": {
				"type": "text"
			},
			"short_icon_url": {
				"type": "text",
				"index": false
			},
			"image_url": {
				"type": "text"
			}
		}
	}
}`

var _ gachifinder.Emitter = &Elasticsearch{}

// Elasticsearch struct
type Elasticsearch struct {
	MajorReleaseNumber  int
	URLs                []string

	// Unexported ...
	client *elastic.Client
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
	e.MajorReleaseNumber = majorReleaseNumber

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
func (e *Elasticsearch) Write(cd <-chan gachifinder.GachiData, done <-chan bool) error {
	// bulkRequest := e.client.Bulk()

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