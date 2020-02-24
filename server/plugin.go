package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/mattermost/mattermost-server/v5/mlog"
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
		ticker := time.NewTicker(time.Second * 10).C
		for {
			select {
			case <-ticker:
				token := p.configuration.Token
				if token != "" {
					mlog.Info(" Bot is running")
					client := NewGoScrumClient("http://192.168.31.56:3000/mattermost/bot", p.configuration.Token)
					workspace, _ := client.GetWorkspaceByToken()
					if workspace == nil {
						return
					}
					if len(workspace.Projects) == 0 {
						return
					}

					for _, project := range workspace.Projects {
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
							createdPost, err := p.CreateBotDMPost(participant.UserID, &post)

							question, _ := client.GetParticipantQuestion(project.ID, participant.ID)

							post = model.Post{
								Message: question.Title,
							}

							createdPost, err = p.CreateBotDMPost(participant.UserID, &post)
							spew.Dump(createdPost)
							spew.Dump(err)
						}
					}
				} else {
					p.API.LogInfo("GoScrum bot API Access Token is required")
				}
			case <-quitChan:
				wg.Done()
				return
			}
		}
	}()

	return nil
}
