package slackbot

import (
	"fmt"

	"github.com/nlopes/slack"

	"../config"
)

var (
	rtm *slack.RTM
)

func Listen() {
	conf := config.Get()

	api := slack.New(conf.SlackToken)
	api.SetDebug(conf.DebugMode)

	// Users List
	/*
		usersByID := make(map[string]*slack.User)
		usersArr, err := api.GetUsers()
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		for _, user := range usersArr {
			usersByID[user.ID] = &user
			//fmt.Printf(" - user: ID=%s Name=%s RealName=%s, Email=%s\n", user.ID, user.Name, user.RealName, user.Profile.Email)
		}*/

	// Open Websocket
	rtm = api.NewRTM()
	go rtm.ManageConnection()

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
