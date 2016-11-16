package slackbot

import (
	"fmt"

	"github.com/nlopes/slack"

	"../links"
)

type Message struct {
	originalMsg     *slack.Msg // Message posted by user
	extendedLinkMsg *slack.Msg // Message with extended links informations
}

func NewMessageFromEvent(eventData *slack.MessageEvent) *Message {
	m := new(Message)
	m.originalMsg = &eventData.Msg
	m.extendedLinkMsg = eventData.SubMessage
	return m
}

func (m *Message) Analyse() bool {

	// Ignore bot messages!
	if len(m.originalMsg.BotID) > 0 {
		return false
	}

	// There is one link?
	links := m.GetLinks()
	if len(links) == 0 {
		return false
	}

	// Save all links!
	for _, link := range links {

		// link already posted?
		duplicates := link.FindDuplicates()
		if len(duplicates) > 0 {
			// TODO auteur + date
			rtm.SendMessage(rtm.NewOutgoingMessage("Pssst! Someone already posted this link!", m.originalMsg.Channel))
			continue
		}

		fmt.Printf("\n\n------------------\n\n")
		rtm.SendMessage(rtm.NewOutgoingMessage("Thank you!", m.originalMsg.Channel))
		link.Save()
	}

	return true
}

func (m *Message) GetLinks() []*links.Link {
	var messagelinks []*links.Link

	// No SubMessage, No Attachment, No Link
	if m.extendedLinkMsg == nil || len(m.extendedLinkMsg.Attachments) == 0 {
		return messagelinks
	}

	// Links Filter
	for _, attachment := range m.extendedLinkMsg.Attachments {
		if len(attachment.TitleLink) > 0 {
			fmt.Printf("- Link: %s, %s, %s\n", attachment.Title, attachment.TitleLink, attachment.Text)
			link := links.NewLink(attachment.TitleLink)
			link.SetTitle(attachment.Title)
			messagelinks = append(messagelinks, link)
		}
	}

	return messagelinks
}
