package main

import (
	"flag"

	"./src/api"
	"./src/config"
	"./src/links"
	"./src/slackbot"
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
