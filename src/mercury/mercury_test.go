// +build !travis

package mercury

import (
	"flag"
	"testing"

	"github.com/jguyomard/slackbot-links/src/config"
)

var (
	testURL                   = "https://ilonet.fr/welcome-to-hugo/"
	testExpectedTitle         = "Welcome to Hugo!"
	testExpectedDatePublished = "2016-06-11T00:00:00.000Z"
	testExpectedImageURL      = ""
	testExpectedTotalPages    = 1
	testExpectedWordCount     = 92
	testExpectedDirection     = "ltr"
)

func init() {
	configFilePtr := flag.String("config-file", "/etc/slackbot-links/config.yaml", "conf file path")
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

	if infos.Direction != testExpectedDirection {
		t.Fatal("mercury.Parse(): incorrect title")
	}

	if infos.WordCount != testExpectedWordCount {
		t.Fatal("mercury.Parse(): incorrect title")
	}
}
