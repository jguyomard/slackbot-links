package links

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"
)

const (
	esIndex = "slackbot-bookmarks6"
	esType  = "links"
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
	es, err = elastic.NewClient()
	if err != nil {
		panic(err)
	}

	// Create Index/Type?
	createESIndexIfNeeded()
}

func createESIndexIfNeeded() bool {

	// L'index n'existe pas déjà ?
	exists, err := es.IndexExists(esIndex).Do()
	if err != nil {
		panic(err)
	}
	if exists {
		fmt.Printf("Index %s already exists.\n", esIndex)
		return false
	}

	// Création de l'index
	fmt.Printf("Create new index: %s...\n", esIndex)
	createIndex, err := es.CreateIndex(esIndex).Do()
	if err != nil {
		panic(err)
	}
	if !createIndex.Acknowledged {
		panic("ELASTIC FATAL ERROR: CreateIndex failed!")
	}

	mapping := `{
			"properties":{
        "url":{
					"type":     "string",
					"index":    "not_analyzed"
				},
				"date_published":{
					"type":"date"
				}
			}
	}`
	putMapping, err := es.PutMapping().
		Index(esIndex).
		Type(esType).
		BodyString(mapping).
		Do()
	if err != nil {
		panic(err)
	}
	if !putMapping.Acknowledged {
		panic("ELASTIC FATAL ERROR: PutMapping failed!")
	}

	return true
}
