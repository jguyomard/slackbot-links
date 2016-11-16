package links

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"gopkg.in/olivere/elastic.v3"
)

func Search(params url.Values) []*Link {

	searchService := es.Search().
		Index(esIndex).
		Type(esType).
		From(getLimitOffset(params)).
		Size(getLimitSize(params)).
		SortWithInfo(getSortInfo(params)).
		Pretty(true)

	query := getQuery(params)
	if query != nil {
		searchService.Query(query)
	}

	searchResult, err := searchService.Do()
	if err != nil {
		fmt.Printf("Links::Search() Errors : %s\n", err)
		return []*Link{}
	}
	return searchResultsToLinks(searchResult)
}

func getQuery(params url.Values) *elastic.BoolQuery {
	search := params.Get("search")
	url := params.Get("url")
	if search == "" && url == "" {
		return nil
	}

	query := elastic.NewBoolQuery()
	if search != "" {
		query.Must(elastic.NewMultiMatchQuery(search, "title^8"))
	}
	if url != "" {
		query.Must(elastic.NewTermQuery("url", url))
	}
	return query
}

func getSortInfo(params url.Values) elastic.SortInfo {
	return elastic.SortInfo{Field: "date_published", Ascending: false}
}

func getLimitOffset(params url.Values) int {
	return getIntParamsOrDefault(params, "offset", 0)
}

func getLimitSize(params url.Values) int {
	return getIntParamsOrDefault(params, "size", 100)
}

func getIntParamsOrDefault(params url.Values, paramKey string, defaultValue int) int {
	val := params.Get(paramKey)
	if val == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return intVal
}

func searchResultsToLinks(searchResult *elastic.SearchResult) []*Link {
	var links []*Link

	var ttyp Link
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if link, ok := item.(Link); ok {
			links = append(links, &link)
			fmt.Printf("url=%s, title=%s\n", link.URL, link.Title)
		}
	}

	return links
}
