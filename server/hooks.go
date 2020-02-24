package main

import (
	"github.com/davecgh/go-spew/spew"
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

	spew.Dump(post)

	// Make sure this is a post sent directly to Surveybot
	channel, appErr := p.API.GetChannel(post.ChannelId)
	if appErr != nil {
		p.API.LogError("Unable to get channel for Surveybot feedback", "err", appErr)
		return
	}

	if !p.IsBotDMChannel(channel) {
		return
	}

	mlog.Info("Done")
	// Make sure this is not a post sent by another bot
	//user, appErr := p.API.GetUser(post.UserId)
	//if appErr != nil {
	//	p.API.LogError("Unable to get sender for Surveybot feedback", "err", appErr)
	//	return
	//}
	//
	//if user.IsBot {
	//	return
	//}

	//// Send the feedback to Segment
	//if err := p.sendFeedback(post.Message, post.UserId, post.CreateAt); err != nil {
	//	p.API.LogError("Failed to send Surveybot feedback to Segment", "err", err.Error())
	//
	//	// Still appear to the end user as if their feedback was actually sent
	//}
	//
	//// Respond to the feedback
	//_, appErr = p.CreateBotDMPost(post.UserId, &model.Post{
	//	Message: feedbackResponseBody,
	//	Type:    "custom_nps_thanks",
	//})
	if appErr != nil {
		p.API.LogError("Failed to respond to Surveybot feedback")
	}
}

//func (p *Plugin) UserHasLoggedIn(c *plugin.Context, user *model.User) {
//	if err := p.checkForDMs(user.Id); err != nil {
//		p.API.LogError("Failed to check for user notifications on login", "user_id", user.Id, "err", err)
//	}
//}