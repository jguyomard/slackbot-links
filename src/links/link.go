package links

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kennygrant/sanitize"
	"github.com/satori/go.uuid"

	"github.com/jguyomard/slackbot-links/src/mercury"
)

type Link struct {
	ID            string     `json:"id"`
	URL           string     `json:"url"`
	Title         string     `json:"title"`
	Author        string     `json:"author"`
	Excerpt       string     `json:"excerpt"`
	DatePublished *time.Time `json:"published_at,omitempty"`
	ImageURL      string     `json:"image_url"`
	Content       string     `json:"content"`
	SharedAt      *time.Time `json:"shared_at"`
	SharedBy      struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"shared_by"`
	SharedOn struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"shared_on"`
}

func NewLink(url string) *Link {
	Init()

	l := new(Link)
	l.URL = url
	now := time.Now()
	l.SharedAt = &now
	return l
}

func (l *Link) SetTitle(title string) {
	l.Title = title
}

func (l *Link) SetExcerpt(excerpt string) {
	l.Excerpt = sanitize.HTML(excerpt)
}

func (l *Link) SetImageURL(imageURL string) {
	l.ImageURL = imageURL
}

func (l *Link) SetSharedAt(date *time.Time) {
	l.SharedAt = date
}

func (l *Link) SetSharedBy(userID, userName string) {
	l.SharedBy.ID = userID
	l.SharedBy.Name = userName
}

func (l *Link) SetSharedOn(channelID, channelName string) {
	l.SharedOn.ID = channelID
	l.SharedOn.Name = channelName
}

// GetID returns current ID or generate new one
func (l *Link) GetID() string {
	if l.ID == "" {
		l.ID = uuid.NewV4().String()
	}
	return l.ID
}

// FindDuplicates checks if this link is already posted?
func (l *Link) FindDuplicates() *SearchResult {
	params := url.Values{}
	params.Set("url", l.URL)
	search, _ := Search(params)
	return search
}

// Save this link to Elastic Search
func (l *Link) Save() bool {

	// Log
	log.WithFields(log.Fields{
		"action": "save",
		"linkID": l.GetID(),
		"link":   l,
	}).Info("New link")

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

func (l *Link) Delete() bool {

	// Log
	log.WithFields(log.Fields{
		"action":    "delete",
		"linkID":    l.GetID(),
		"linkURL":   l.URL,
		"linkTitle": l.Title,
	}).Info("Delete link")

	// Delete from Elastic Search
	_, err := es.Delete().
		Index(esIndex).
		Type(esType).
		Id(l.GetID()).
		Do()

	if err != nil {
		fmt.Printf("ES Error: %s\n", err)
		return false
	}

	return true
}

func (l *Link) enrichMetasFromContent() {
	infos, err := mercury.ParseURL(l.URL)
	if err != nil {
		return
	}

	if len(infos.Title) > 0 {
		l.SetTitle(infos.Title)
	}
	if len(infos.Excerpt) > 0 {
		l.SetExcerpt(infos.Excerpt)
	}
	if len(infos.ImageURL) > 0 {
		l.SetImageURL(infos.ImageURL)
	}

	l.Author = infos.Author
	l.Content = sanitize.HTML(infos.Content)

	datePublished, err := time.Parse("2006-01-02T15:04:05.000Z", infos.DatePublished)
	if err == nil {
		l.DatePublished = &datePublished
	}
}
