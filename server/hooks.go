package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	mlog.Info("Message has been posted")
	// Make sure that Surveybot doesn't respond to itself
	if post.UserId == p.botUserID {
		return
	}

	mlog.Info("Check system message")

	// Or to system messages
	if post.IsSystemMessage() {
		return
	}

	mlog.Info("Check user channel")

	// Make sure this is a post sent directly to GoScrum
	channel, appErr := p.API.GetChannel(post.ChannelId)
	if appErr != nil {
		p.API.LogError("Unable to get channel for Surveybot feedback", "err", appErr)
		return
	}

	if !p.IsBotDMChannel(channel) {
		return
	}

	mlog.Info("Done")
	if p.configuration.Token != "" {
		p.goScrumClient = NewGoScrumClient(fmt.Sprintf("%s/mattermost/bot", p.configuration.URL), p.configuration.Token)
		// TODO -- check for failure cases.
		message, _ := p.goScrumClient.UserInteraction(post.UserId, post.Message)

		if message.MessageType == QuestionMessage {
			if message.Question.Title != "" {
				newPost := model.Post{
					Message: message.Question.Title,
				}

				createdPost, err := p.CreateBotDMPost(post.UserId, &newPost)
				if err != nil {
					fmt.Println("Error:", err.Error())
				}
				if createdPost != nil {
					p.goScrumClient.UpdateAnswerPost(message.ParticipantID, message.Question.ID, createdPost.Id)
				}
			}
		}
	}

}

//func (p *Plugin) UserHasLoggedIn(c *plugin.Context, user *model.User) {
//	if err := p.checkForDMs(user.Id); err != nil {
//		p.API.LogError("Failed to check for user notifications on login", "user_id", user.Id, "err", err)
//	}
//}
