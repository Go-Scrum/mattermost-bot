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
					_ = p.goScrumClient.AddUserActivity(UserActivity{
						UserId:        createdPost.UserId,
						ChannelID:     createdPost.ChannelId,
						ProjectID:     message.Question.ProjectId,
						ParticipantID: message.ParticipantID,
						QuestionID:    message.Question.ID,
						BotPostId:     p.botUserID,
						ActivityType:  UserQuestionActivity,
					})
				}
			}
		}

		if message.MessageType == StandupMessage {
			newPost := &model.Post{}
			model.ParseSlackAttachment(newPost, message.Attachments)
			_, err := p.CreateBotDMPost(post.UserId, newPost)
			if err != nil {
				fmt.Println("Error:", err.Error())
			}
			//if createdPost != nil {
			//	_ = p.goScrumClient.AddUserActivity(UserActivity{
			//		UserId:        createdPost.UserId,
			//		ChannelID:     createdPost.ChannelId,
			//		ProjectID:     message.Question.ProjectId,
			//		ParticipantID: message.ParticipantID,
			//		QuestionID:    message.Question.ID,
			//		BotPostId:     p.botUserID,
			//		ActivityType:  UserStandupActivity,
			//	})
			//}
		}

		if message.MessageType == ReportMessage {
			newPost := &model.Post{
				Message:   message.Content,
				ChannelId: message.ChannelId,
				UserId:    p.botUserID,
			}
			model.ParseSlackAttachment(newPost, message.Attachments)

			createdPost, err := p.CreateChannelPost(newPost)
			if err != nil {
				fmt.Println("Error:", err.Error())
			}
			if createdPost != nil {
				_ = p.goScrumClient.AddUserActivity(UserActivity{
					UserId:        createdPost.UserId,
					ChannelID:     createdPost.ChannelId,
					ParticipantID: message.ParticipantID,
					QuestionID:    message.Question.ID,
					BotPostId:     p.botUserID,
					ActivityType:  UserReportActivity,
				})
			}
		}
	}

}

//func (p *Plugin) UserHasLoggedIn(c *plugin.Context, user *model.User) {
//	if err := p.checkForDMs(user.Id); err != nil {
//		p.API.LogError("Failed to check for user notifications on login", "user_id", user.Id, "err", err)
//	}
//}
