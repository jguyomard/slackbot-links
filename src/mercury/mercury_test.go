package mercury

import (
	"flag"
	"testing"

	"../config"
)

var (
	testURL                   = "http://marcelpociot.com/blog/2016-10-19-write-your-own-slack-bot-using-laravel"
	testExpectedTitle         = "Write a Slack bot using Laravel and PHP"
	testExpectedDatePublished = "2016-10-19T00:00:00.000Z"
	testExpectedImageURL      = "http://marcelpociot.com/user/pages/blog/2016-10-19-write-your-own-slack-bot-using-laravel/Screen%20Shot%202016-10-19%20at%2010.36.15.png"
	testExpectedTotalPages    = 1
)

func init() {
	configFilePtr := flag.String("config-file", "/etc/slack-bookmarks/config.yaml", "conf file path")
	flag.Parse()
	config.SetFilePath(*configFilePtr)
}

func TestParse(t *testing.T) {
	infos, err := ParseURL(testURL)

	if err != nil {
		t.Fatal(err)
	}

	if infos.Title != testExpectedTitle {
		t.Fatal("mercury.Parse(): incorrect Title")
	}

	if infos.DatePublished != testExpectedDatePublished {
		t.Fatal("mercury.Parse(): incorrect DatePublished")
	}

	if infos.ImageURL != testExpectedImageURL {
		t.Fatal("mercury.Parse(): incorrect ImageURL")
	}

	if infos.TotalPages != testExpectedTotalPages {
		t.Fatal("mercury.Parse(): incorrect title")
	}
}
