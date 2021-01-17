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

var _ gachifinder.Emitter = &Elasticsearch{}

// Elasticsearch struct
type Elasticsearch struct {
	MajorReleaseNumber  int
	URLs                []string

	Client *elastic.Client
}

// Connect to Elasticsearch & Create index.
func (e *Elasticsearch) Connect() error {
	_, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
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

	e.Client = client
	e.MajorReleaseNumber = majorReleaseNumber

	return nil
}

// Close to release Elasticsearch Client.
func (e *Elasticsearch) Close() {
	e.Client = nil
}

// Write the data into the Elasticsearch.
func (e *Elasticsearch) Write() {

}