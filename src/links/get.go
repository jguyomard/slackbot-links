package links

import (
	"encoding/json"
	"fmt"
)

// Get allows to get a link from Elastic Search based on its id
func Get(id string) (*Link, bool) {

	// Get by id from Elastic Search
	getResult, err := es.Get().
		Index(esIndex).
		Type(esType).
		Id(id).
		Do()
	if err != nil {
		return nil, false
	}

	// Json Decode
	var link Link
	linkJSON, _ := getResult.Source.MarshalJSON()
	if err = json.Unmarshal(linkJSON, &link); err != nil {
		fmt.Printf("err=%s\n", err)
		return nil, false
	}

	return &link, true
}
