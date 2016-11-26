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
		if duplicates.GetTotal() > 0 {
			duplicateLink := duplicates.GetLinks()[0]
			duplicateAuthor := duplicateLink.SharedBy.Name
			if duplicateAuthor == "" {
				duplicateAuthor = "Someone"
			}
			// TODO add duplicate date?
			duplicateMsg := fmt.Sprintf("Pssst! %s already posted this link!", duplicateAuthor)
			rtm.SendMessage(rtm.NewOutgoingMessage(duplicateMsg, m.originalMsg.Channel))
			continue
		}

		//rtm.SendMessage(rtm.NewOutgoingMessage("Link saved, Thank you!", m.originalMsg.Channel))
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
			link := links.NewLink(attachment.TitleLink)
			link.SetTitle(attachment.Title)
			link.SetSharedBy(m.extendedLinkMsg.User, getUserName(m.extendedLinkMsg.User))
			link.SetSharedOn(m.originalMsg.Channel, getChannelName(m.originalMsg.Channel))
			messagelinks = append(messagelinks, link)
		}
	}

	return messagelinks
}

func getUserName(userID string) string {
	userInfos, err := rtm.GetUserInfo(userID)
	if err != nil {
		return ""
	}
	return userInfos.Name
}

func getChannelName(channelID string) string {
	channelInfos, err := rtm.GetChannelInfo(channelID)
	if err != nil {
		return ""
	}
	return channelInfos.Name
}
