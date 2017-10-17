package main

import (
	"flag"
	"os"

	"github.com/jguyomard/slackbot-links/src/api"
	"github.com/jguyomard/slackbot-links/src/config"
	"github.com/jguyomard/slackbot-links/src/links"
	"github.com/jguyomard/slackbot-links/src/slackbot"
)

func main() {

	// Read command options
	configFilePtr := flag.String("config-file", "/etc/slackbot-links/config.yaml", "conf file path")
	linksFilePtr := flag.String("links-file", "", "links file path")
	flag.Parse()

	// Set config filepath
	config.SetFilePath(*configFilePtr)

	// Connect to ES
	links.Init()

	// Commands
	if flag.Arg(0) == "restore" {
		if *linksFilePtr == "" {
			panic("--links-file is mandadory.")
		}
		links.Restore(*linksFilePtr)
		os.Exit(0)
	}

	// Listen for new Slack Messages
	go slackbot.Listen()

	// API
	api.Listen()

}
