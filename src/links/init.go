package links

import (
	"context"
	"fmt"
	"os"

	"github.com/jguyomard/slackbot-links/src/config"

	log "github.com/Sirupsen/logrus"
	"github.com/olivere/elastic"
)

const (
	esIndex    = "slackbot-links"
	esType     = "links"
	esAnalyzer = "french"
)

var (
	es *elastic.Client
)

func Init() {
	var err error

	if es != nil {
		return
	}

	// Elastic Search Connection
	es, err = elastic.NewClient(
		elastic.SetURL(config.Get().ElasticSearchURLS...),
		elastic.SetSniff(!config.Get().ElasticSearchDisableSniffing),
	)
	if err != nil {
		panic(err)
	}

	// Logging
	log.SetFormatter(&log.JSONFormatter{})
	logFileName := fmt.Sprintf("%s/%s", config.Get().LogsDir, "links.log")
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)

	// Create Index/Type?
	createESIndexIfNeeded()
}

func createESIndexIfNeeded() bool {

	ctx := context.Background()

	// Index already exists?
	exists, err := es.IndexExists(esIndex).Do(ctx)
	if err != nil {
		panic(err)
	}
	if exists {
		return false
	}

	// Create new index
	fmt.Printf("Create new index: %s... ", esIndex)
	createIndex, err := es.CreateIndex(esIndex).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !createIndex.Acknowledged {
		panic("ELASTIC FATAL ERROR: CreateIndex failed!")
	}

	mapping := `{
		"properties":{
			"url":{
				"type":  "keyword",
				"index": true
			},
			"title":{
				"type": "text",
				"analyzer": "` + esAnalyzer + `"
			},
			"author":{
				"type": "text"
			},
			"excerpt":{
				"type": "text",
				"analyzer": "` + esAnalyzer + `"
			},
			"published_at":{
				"type":"date"
			},
			"image_url":{
				"type":  "keyword",
				"index": true
			},
			"content":{
				"type": "text",
				"analyzer": "` + esAnalyzer + `"
			},
			"shared_at":{
				"type":"date"
			}
		}
	}`
	putMapping, err := es.PutMapping().
		Index(esIndex).
		Type(esType).
		BodyString(mapping).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	if !putMapping.Acknowledged {
		panic("ELASTIC FATAL ERROR: PutMapping failed!")
	}

	fmt.Println("OK")
	return true
}
