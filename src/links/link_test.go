package links

import (
	"flag"
	"testing"

	"github.com/jguyomard/slackbot-links/src/config"
)

var (
	testURL   = "http://marcelpociot.com/blog/2016-10-19-write-your-own-slack-bot-using-laravel"
	testTitle = "a_title"
	testDesc  = "a_desc"
)

func init() {
	configFilePtr := flag.String("config-file", "/etc/slackbot-links/config.yaml", "conf file path")
	flag.Parse()
	config.SetFilePath(*configFilePtr)
}

func TestSaveLink(t *testing.T) {
	link := NewLink(testURL)

	link.SetTitle(testTitle)
	if link.Title != testTitle {
		t.Fatal("link.SetTitle() error")
	}

	saved := link.Save()
	if !saved {
		t.Fatal("link.Save() error")
	}
}

func TestFindDuplicate(t *testing.T) {
	link := NewLink(testURL)
	duplicates := link.FindDuplicates()

	if duplicates.GetTotal() == 0 {
		t.Fatal("link.FindDuplicates() error")
	}
}
