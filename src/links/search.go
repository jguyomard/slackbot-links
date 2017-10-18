package links

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"gopkg.in/olivere/elastic.v5"
)

// SearchResult contains links that match the query, aggregations (stats) and total hits
type SearchResult struct {
	links         []*Link
	total         int
	perPage       int
	currentOffset int
	stats         struct {
		sharedBy map[string]int64
		sharedAt map[string]int64
	}
}

// Search allows to get multiple links from Elastic Search, that match the query
func Search(params url.Values) (*SearchResult, error) {

	ctx := context.Background()

	searchService := es.Search().
		Index(esIndex).
		Type(esType).
		From(getLimitOffset(params)).
		Size(getLimitSize(params)).
		SortBy(getSortInfo(params)...).
		Aggregation("shared_by", elastic.NewTermsAggregation().Field("shared_by.name.keyword").Size(50)).
		Aggregation("shared_at", elastic.NewDateHistogramAggregation().Field("shared_at").Interval("month")).
		Pretty(true)

	query := getQuery(params)
	if query != nil {
		searchService.Query(query)
	}

	elasticResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Printf("Links::Search() Errors : %s\n", err)
		return nil, fmt.Errorf("Invalid Seach Query")
	}

	res := elasticResultsToLinksResult(elasticResult)
	res.perPage = getLimitSize(params)
	res.currentOffset = getLimitOffset(params)
	return res, nil
}

// GetLinks returns current page links
func (r *SearchResult) GetLinks() []*Link {
	return r.links
}

// GetTotal returns total links count (total hits), for this SearchResult
func (r *SearchResult) GetTotal() int {
	return r.total
}

// GetCursor returns pagination (with cursors)
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

// GetStats return stats (aggregations), for this SearchResult
func (r *SearchResult) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"shared_by": r.stats.sharedBy,
		"shared_at": r.stats.sharedAt,
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

	// Links
	var ttyp Link
	for _, item := range elasticResult.Each(reflect.TypeOf(ttyp)) {
		if link, ok := item.(Link); ok {
			res.links = append(res.links, &link)
		}
	}

	// Stats: Deserialize aggregations
	if agg, found := elasticResult.Aggregations.Terms("shared_by"); found {
		res.stats.sharedBy = make(map[string]int64)
		for _, bucket := range agg.Buckets {
			res.stats.sharedBy[bucket.Key.(string)] = bucket.DocCount
		}
	}
	if agg, found := elasticResult.Aggregations.Histogram("shared_at"); found {
		res.stats.sharedAt = make(map[string]int64)
		for _, bucket := range agg.Buckets {
			res.stats.sharedAt[(*bucket.KeyAsString)[:7]] = bucket.DocCount
		}
	}

	// Total
	res.total = int(elasticResult.TotalHits())

	return res
}
