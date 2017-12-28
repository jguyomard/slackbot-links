package links

import (
	"net/url"
	"testing"
)

func TestSearchParamsParse(t *testing.T) {
	params := url.Values{}
	params.Add("offset", "13")
	params.Add("limit", "37")
	params.Add("invalid", "foo")
	if getIntParamsOrDefault(params, "invalid", 42) != 42 {
		t.Fatal("TestSearchParamsParse: getIntParamsOrDefault() error")
	}
	if getIntParamsOrDefault(params, "not-set", 42) != 42 {
		t.Fatal("TestSearchParamsParse: getIntParamsOrDefault() error")
	}
	if getLimitOffset(params) != 13 {
		t.Fatal("TestSearchParamsParse: getLimitOffset() error")
	}
	if getLimitSize(params) != 37 {
		t.Fatal("TestSearchParamsParse: getLimitSize() error")
	}
}

func TestSearchLink_byTitle(t *testing.T) {
	link := NewLink(testURL)
	link.SetTitle(testTitle)
	link.Save()

	params := url.Values{}
	params.Add("search", testTitleSearch)
	search, err := Search(params)
	link.Delete()

	if err != nil {
		t.Fatalf("TestSearchLink_byTitle: Search() error, %s", err)
	}
	if search.GetTotal() == 0 {
		t.Fatal("TestSearchLink_byTitle: Search() error, no result found")
	}
	if len(search.GetLinks()) == 0 {
		t.Fatal("TestSearchLink_byTitle: Search() error, no result returned")
	}
}

func TestSearchLink_byExcerpt(t *testing.T) {
	link := NewLink(testURL)
	link.SetExcerpt(testExcerpt)
	link.Save()

	params := url.Values{}
	params.Add("search", testExcerptSearch)
	search, err := Search(params)
	link.Delete()

	if err != nil {
		t.Fatalf("TestSearchLink_byExcerpt: Search() error, %s", err)
	}
	if search.GetTotal() == 0 {
		t.Fatal("TestSearchLink_byExcerpt: Search() error, no result found")
	}
}

func TestSearchLink_bySharedBy(t *testing.T) {
	link := NewLink(testURL)
	link.SetSharedBy("1", testSharedBy)
	link.Save()

	params := url.Values{}
	params.Add("sharedby", testSharedBy)
	search, err := Search(params)
	link.Delete()

	if err != nil {
		t.Fatalf("TestSearchLink_bySharedBy: Search() error, %s", err)
	}
	if search.GetTotal() == 0 {
		t.Fatal("TestSearchLink_bySharedBy: Search() error, no result found")
	}
}

func TestSearchLink_byURL(t *testing.T) {
	link := NewLink(testURL)
	link.Save()

	params := url.Values{}
	params.Add("url", testURL)
	search, err := Search(params)
	link.Delete()

	if err != nil {
		t.Fatalf("TestSearchLink_byURL: Search() error, %s", err)
	}
	if search.GetTotal() == 0 {
		t.Fatal("TestSearchLink_byURL: Search() error, no result found")
	}
}

func TestEmptyQuery(t *testing.T) {
	params := url.Values{}
	query := getQuery(params)
	if query != nil {
		t.Fatal("TestEmptyQuery: Search() error, nil expected")
	}
}

func TestGetCursor(t *testing.T) {
	link := NewLink(testURL)
	link.SetSharedBy("1", testSharedBy)
	link.Save()

	params := url.Values{}
	params.Add("sharedby", testSharedBy)
	search, err := Search(params)
	link.Delete()

	if err != nil {
		t.Fatalf("TestGetCursor: Search() error, %s", err)
	}

	cursor := search.GetCursor()
	if cursor["previous"] != nil {
		t.Fatal("TestGetCursor: GetCursor()[previous] != nil")
	}
	if cursor["next"] != nil {
		t.Fatal("TestGetCursor: GetCursor()[next] != nil")
	}
}

func TestGetStats(t *testing.T) {
	link := NewLink(testURL)
	link.SetSharedBy("1", testSharedBy)
	link.Save()

	params := url.Values{}
	params.Add("sharedby", testSharedBy)
	search, err := Search(params)
	link.Delete()

	if err != nil {
		t.Fatalf("TestGetStats: Search() error, %s", err)
	}

	stats, ok := search.GetStats()["shared_by"].(map[string]int64)
	if !ok {
		t.Fatal("TestGetStats: GetStats() error, shared_by isn't a map[string]int64")
	}
	if stats[testSharedBy] != 1 {
		t.Fatalf("TestGetStats: GetStats() error, invalid aggregation for %s, %d instead of 1", testSharedBy, stats[testSharedBy])
	}
}
