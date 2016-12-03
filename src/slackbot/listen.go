package slackbot

import (
	"fmt"

	"github.com/nlopes/slack"

	"github.com/jguyomard/slackbot-links/src/config"
)

var (
	rtm *slack.RTM
)

// Listen connect to Slack through websocket
func Listen() {
	conf := config.Get()

	api := slack.New(conf.SlackToken)
	api.SetDebug(conf.DebugMode)

	// Open Websocket
	rtm = api.NewRTM()
	go rtm.ManageConnection()

	// Check Auth
	_, err := rtm.AuthTest()
	if err != nil {
		panic("Unable to connect to Slack ; invalid token?")
	}

	for {
		slackEvent := <-rtm.IncomingEvents

		switch eventData := slackEvent.Data.(type) {

		case *slack.MessageEvent:
			if conf.DebugMode {
				fmt.Printf("Incoming Message Event: %+v\n\n", eventData)
				fmt.Printf("Incoming SubMessage Event: %+v\n\n", eventData.SubMessage)
			}
			NewMessageFromEvent(eventData).Analyse()

		default:
			//fmt.Printf("Unknown Event Received: %+v\n", eventData)

		}
	}

}
