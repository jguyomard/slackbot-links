package links

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/satori/go.uuid"

	"../mercury"
)

type Link struct {
	ID            string     `json:"id"`
	URL           string     `json:"url"`
	Title         string     `json:"title"`
	Author        string     `json:"author"`
	Excerpt       string     `json:"excerpt"`
	DatePublished *time.Time `json:"date_published,omitempty"`
	ImageURL      string     `json:"image_url"`
	Content       string     `json:"content"`
}

func NewLink(url string) *Link {
	Init()

	l := new(Link)
	l.URL = url
	return l
}

func (l *Link) SetTitle(title string) {
	l.Title = title
}

func (l *Link) GetID() string {
	if l.ID == "" {
		l.ID = uuid.NewV4().String()
	}
	return l.ID
}

// FindDuplicates checks if this link is already posted?
func (l *Link) FindDuplicates() []*Link {
	params := url.Values{}
	params.Set("url", l.URL)
	return Search(params)
}

// Save this link to Elastic Search
func (l *Link) Save() bool {

	// Enrich this link
	l.enrichMetasFromContent()

	// Save to Elastic Search
	_, err := es.Index().
		Index(esIndex).
		Type(esType).
		Id(l.GetID()).
		BodyJson(l).
		Refresh(true).
		Do()

	if err != nil {
		resJSON, _ := json.Marshal(l)
		fmt.Printf("ES Error: %s\n  -> link=%+v\n  -> linkJSON=%s\n", err, l, string(resJSON))
		return false
	}

	return true
}

func (l *Link) enrichMetasFromContent() {
	infos, _ := mercury.ParseURL(l.URL)

	l.Title = infos.Title
	l.Author = infos.Author
	l.Excerpt = infos.Excerpt
	l.ImageURL = infos.ImageURL
	l.Content = infos.Content

	datePublished, err := time.Parse("2006-01-02T15:04:05.000Z", infos.DatePublished)
	if err == nil {
		l.DatePublished = &datePublished
	}
}
