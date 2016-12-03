package mercury

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jguyomard/slackbot-links/src/config"
)

const (
	apiParseURI = "https://mercury.postlight.com/parser"
)

type ParseInfos struct {
	Title         string
	Content       string
	Excerpt       string
	Author        string
	DatePublished string `json:"date_published"`
	ImageURL      string `json:"lead_image_url"`
	URL           string
	Domain        string
	WordCount     int `json:"word_count"`
	Direction     string
	TotalPages    int    `json:"total_pages"`
	ErrorMessage  string `json:"errorMessage"`
}

// ParseURL transforms web pages into clean text using Mercury Parser
func ParseURL(contentURL string) (*ParseInfos, error) {
	apiKey := config.Get().MercuryAPIKey
	if apiKey == "" {
		return nil, errors.New("mercury : No apiKey, no parsing")
	}

	// API Request
	client := &http.Client{}

	params := url.Values{}
	params.Add("url", contentURL)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", apiParseURI, params.Encode()), nil)
	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("Accept", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	// Response Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// Json Decode
	var infos ParseInfos
	if err := json.Unmarshal(body, &infos); err != nil {
		return nil, err
	}
	if r.StatusCode >= 400 {
		return &infos, fmt.Errorf("Mercury StatusCode Error : %d\n", r.StatusCode)
	}

	return &infos, nil
}
