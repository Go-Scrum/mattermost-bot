package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/1set/cronrange"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	botUsername    = "GoScrum"
	botDisplayName = "GoScrum"
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
}

// OnActivate ensure the bot account exists
func (p *Plugin) OnActivate() error {
	systemBot := &model.Bot{
		Username:    botUsername,
		DisplayName: botDisplayName,
		Description: botDescription,
	}
	botUserID, appErr := p.Helpers.EnsureBot(systemBot)
	if appErr != nil {
		return errors.Wrap(appErr, "failed to ensure systemBot user")
	}
	p.botUserID = botUserID

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
								message := fmt.Sprintf("Hello %s :wave: It's time for **%s** in # %s\n Please share what you've been working on",
									participant.RealName,
									project.Name,
									project.ChannelName,
								)
								post := model.Post{
									// TODO - format the name based on first name and lastName
									Message: message,
								}
								_, _ = p.CreateBotDMPost(participant.UserID, &post) // TODO

								question, _ := p.goScrumClient.GetQuestionDetails(project.Questions[0].ID)

								fmt.Println("Question", question.ID)
								if question.Title != "" {
									post = model.Post{
										Message: question.Title,
									}

									createdPost, _ := p.CreateBotDMPost(participant.UserID, &post)
									p.goScrumClient.UpdateAnswerPost(participant.ID, question.ID, createdPost.Id)
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
