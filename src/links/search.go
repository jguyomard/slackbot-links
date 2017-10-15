package links

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"gopkg.in/olivere/elastic.v3"
)

type SearchResult struct {
	links         []*Link
	total         int
	perPage       int
	currentOffset int
}

// Search allows to get multiple links from Elastic Search, that match the query
func Search(params url.Values) (*SearchResult, error) {

	searchService := es.Search().
		Index(esIndex).
		Type(esType).
		From(getLimitOffset(params)).
		Size(getLimitSize(params)).
		SortBy(getSortInfo(params)...).
		Pretty(true)

	query := getQuery(params)
	if query != nil {
		searchService.Query(query)
	}

	elasticResult, err := searchService.Do()
	if err != nil {
		fmt.Printf("Links::Search() Errors : %s\n", err)
		return nil, fmt.Errorf("Invalid Seach Query")
	}

	res := elasticResultsToLinksResult(elasticResult)
	res.perPage = getLimitSize(params)
	res.currentOffset = getLimitOffset(params)
	return res, nil
}

func (r *SearchResult) GetLinks() []*Link {
	return r.links
}

func (r *SearchResult) GetTotal() int {
	return r.total
}

func (r *SearchResult) GetCursor() map[string]interface{} {

	var previous interface{}
	if r.currentOffset-r.perPage >= 0 {
		previous = r.currentOffset - r.perPage
	} else {
		previous = nil
	}

	var next interface{}
	if r.currentOffset+r.perPage < r.total {
		next = r.currentOffset + r.perPage
	} else {
		next = nil
	}

	return map[string]interface{}{
		"previous": previous,
		"current":  r.currentOffset,
		"next":     next,
		"per_page": r.perPage,
	}
}

func getQuery(params url.Values) *elastic.BoolQuery {
	search := params.Get("search")
	url := params.Get("url")
	sharedby := params.Get("sharedby")
	if search == "" && url == "" && sharedby == "" {
		return nil
	}

	query := elastic.NewBoolQuery()
	if search != "" {
		matchQuery := elastic.NewSimpleQueryStringQuery(search).
			DefaultOperator("and").
			Field("shared_by.name^4").
			Field("shared_on.name^4").
			Field("title^3").
			Field("url^3").
			Field("author^2").
			Field("excerpt^2").
			Field("content^1")
		query.Must(matchQuery)
	}
	if url != "" {
		query.Must(elastic.NewTermQuery("url", url))
	}
	if sharedby != "" {
		query.Must(elastic.NewTermQuery("shared_by.name", sharedby))
	}
	return query
}

func getSortInfo(params url.Values) []elastic.Sorter {
	var sorters []elastic.Sorter
	if params.Get("search") != "" {
		sorters = append(sorters, elastic.NewScoreSort())
	}
	sorters = append(sorters, elastic.SortInfo{Field: "shared_at", Ascending: false})
	return sorters
}

func getLimitOffset(params url.Values) int {
	return getIntParamsOrDefault(params, "offset", 0)
}

func getLimitSize(params url.Values) int {
	return getIntParamsOrDefault(params, "limit", 50)
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

func elasticResultsToLinksResult(elasticResult *elastic.SearchResult) *SearchResult {
	res := new(SearchResult)

	var ttyp Link
	for _, item := range elasticResult.Each(reflect.TypeOf(ttyp)) {
		if link, ok := item.(Link); ok {
			res.links = append(res.links, &link)
		}
	}

	res.total = int(elasticResult.TotalHits())

	return res
}
