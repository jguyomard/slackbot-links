package slackbot

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/nlopes/slack"

	"github.com/jguyomard/slackbot-links/src/links"
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

	// Ignore some messages
	if !m.isValidMessage() {
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

			// if the same author has posted this link recently, it's OK
			if duplicateLink.SharedBy.ID != link.SharedBy.ID || time.Since(*duplicateLink.SharedAt) > 6*time.Hour {
				duplicateAuthor := "Someone"
				if duplicateLink.SharedBy.ID == link.SharedBy.ID {
					duplicateAuthor = "You"
				} else if duplicateLink.SharedBy.Name != "" {
					duplicateAuthor = duplicateLink.SharedBy.Name
				}

				duplicateHumanDate := "a long while ago"
				if duplicateLink.SharedAt != nil {
					duplicateHumanDate = humanize.Time(*duplicateLink.SharedAt)
				}

				duplicateMsg := fmt.Sprintf("@%s Pssst! %s already posted this link %s!", link.SharedBy.Name, duplicateAuthor, duplicateHumanDate)
				rtm.SendMessage(rtm.NewOutgoingMessage(duplicateMsg, m.originalMsg.Channel))
			}

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
			link.SetExcerpt(attachment.Text)
			link.SetImageURL(attachment.ImageURL)
			link.SetSharedBy(m.extendedLinkMsg.User, getUserName(m.extendedLinkMsg.User))
			link.SetSharedOn(m.originalMsg.Channel, getChannelName(m.originalMsg.Channel))
			messagelinks = append(messagelinks, link)
		}
	}

	return messagelinks
}

func (m *Message) isValidMessage() bool {

	// Ignore bot messages!
	if len(m.originalMsg.BotID) > 0 {
		return false
	}

	// Ignore replies, etc
	if m.originalMsg.SubType != "message_changed" {
		return false
	}

	// Ignore messages without attachments
	if m.extendedLinkMsg == nil || len(m.extendedLinkMsg.Attachments) == 0 {
		return false
	}

	// This is probably an edit
	if m.extendedLinkMsg.Edited != nil {
		return false
	}

	return true
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
