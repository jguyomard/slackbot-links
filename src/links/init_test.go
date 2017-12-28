package links

import (
	"flag"

	"github.com/jguyomard/slackbot-links/src/config"
)

const (
	testURL           = "https://ilonet.fr/welcome-to-hugo/"
	testTitle         = "Welcome to Hugo!"
	testTitleSearch   = "Hugo"
	testExcerpt       = "ce blog est maintenant motorisé par le générateur de site"
	testExcerptSearch = "blog"
	testSharedAt      = "2017-01-01 13:37:00"
	testImageURL      = "https://test.img"
	testSharedBy      = "julien"
	testSharedOn      = "links"
)

var (
	linkID string
)

func init() {
	configFilePtr := flag.String("config-file", "/etc/slackbot-links/config.yaml", "conf file path")
	flag.Parse()
	config.SetFilePath(*configFilePtr)
}
