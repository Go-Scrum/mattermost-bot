package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/1set/cronrange"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	botUsername    = "goscrum"
	botDisplayName = "GoScrum Bot"
	botDescription = "A bot account created by the GoScrum plugin."
)

// Plugin represents the welcome bot plugin
type Plugin struct {
	plugin.MattermostPlugin

	// botUserID of the created bot account.
	botUserID string

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	// GoScrumClient
	goScrumClient GoScrumClient

	router *mux.Router
}

// OnActivate ensure the bot account exists
func (p *Plugin) OnActivate() error {

	p.router = p.InitAPI()

	systemBot := &model.Bot{
		Username:    botUsername,
		DisplayName: botDisplayName,
		Description: botDescription,
	}
	botUserID, appErr := p.Helpers.EnsureBot(systemBot)
	fmt.Println("Bot creation")
	if appErr != nil {
		return errors.Wrap(appErr, "failed to ensure systemBot user")
	}
	p.botUserID = botUserID

	fmt.Println("Bot creation", p.botUserID)

	p.API.RegisterCommand(getCommand())

	quitChan := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		ticker := time.NewTicker(time.Second * 60).C
		for {
			select {
			case <-ticker:
				p.StartBot()
			case <-quitChan:
				wg.Done()
				return
			}
		}
	}()

	return nil
}

func (p *Plugin) StartBot() {
	token := p.configuration.Token
	if token != "" {
		fmt.Println(" Bot is running", time.Now())
		p.goScrumClient = NewGoScrumClient(fmt.Sprintf("%s/mattermost/bot", p.configuration.URL), p.configuration.Token)
		workspace, _ := p.goScrumClient.GetWorkspaceByToken()
		if workspace != nil {
			if workspace.Projects != nil || len(workspace.Projects) > 0 {
				for _, project := range workspace.Projects {
					if cr, err := cronrange.ParseString(project.ReportingTime); err == nil {
						current := time.Now()
						if cr.IsWithin(current) {
							for _, participant := range project.Participants {
								message := fmt.Sprintf("Hello @%s :wave: It's time for **%s** in # %s\n Please share what you've been working on",
									participant.RealName,
									project.Name,
									project.ChannelName,
								)
								post := model.Post{
									// TODO - format the name based on first name and lastName
									Message: message,
								}
								createdPost, _ := p.CreateBotDMPost(participant.UserID, &post)
								// TODO check for errors
								_ = p.goScrumClient.AddUserActivity(UserActivity{
									UserId:        createdPost.UserId,
									ChannelID:     createdPost.ChannelId,
									ProjectID:     project.ID,
									ParticipantID: participant.ID,
									QuestionID:    "",
									BotPostId:     p.botUserID,
									ActivityType:  UserGreetingActivity,
								})

								question, _ := p.goScrumClient.GetQuestionDetails(project.Questions[0].ID)

								fmt.Println("Question", question.ID)
								if question.Title != "" {
									post = model.Post{
										Message: question.Title,
									}

									createdPost, _ = p.CreateBotDMPost(participant.UserID, &post)
									_ = p.goScrumClient.AddUserActivity(UserActivity{
										UserId:        createdPost.UserId,
										ChannelID:     createdPost.ChannelId,
										ProjectID:     project.ID,
										ParticipantID: participant.ID,
										QuestionID:    question.ID,
										BotPostId:     p.botUserID,
										ActivityType:  UserQuestionActivity,
									})
								}
							}
						}
					} else {
						fmt.Println("got parse err:", err)
					}
				}
			}
		}
	} else {
		p.API.LogInfo("GoScrum bot API Access Token is required")
	}
}
