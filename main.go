package main

import (
	"flag"

	"github.com/jguyomard/slackbot-links/src/api"
	"github.com/jguyomard/slackbot-links/src/config"
	"github.com/jguyomard/slackbot-links/src/links"
	"github.com/jguyomard/slackbot-links/src/slackbot"
)

func main() {

	// Read command options
	configFilePtr := flag.String("config-file", "/etc/slackbot-links/config.yaml", "conf file path")
	flag.Parse()

	// Set config filepath
	config.SetFilePath(*configFilePtr)

	// Connect to ES
	links.Init()

	// Listen for new Slack Messages
	go slackbot.Listen()

	// API
	api.Listen()

}
