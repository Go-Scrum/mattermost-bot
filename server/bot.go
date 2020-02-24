package main

import (
	"sync"
	"time"

	"github.com/mattermost/mattermost-server/v5/mlog"
)

// Bot struct used for storing and communicating with slack api
type Bot struct {
	botId    string
	conf     *configuration
	quitChan chan struct{}
}

//New creates new Bot instance
func New(botId string, config *configuration) *Bot {
	bot := &Bot{
		conf:  config,
		botId: botId,
	}
	bot.quitChan = make(chan struct{})
	return bot
}

//Start updates Users list and launches notifications
func (bot *Bot) Start() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		ticker := time.NewTicker(time.Second * 10).C
		for {
			select {
			case <-ticker:
				if bot.conf.Token != "" {
					mlog.Info(" Bot is running")
					client := NewGoScrumClient("http://192.168.31.56:3000/mattermost/bot", bot.conf.Token)
					workspace, _ := client.GetWorkspaceByToken()
					if workspace != nil {
						// send greeting message to users
					}

				}
			case <-bot.quitChan:
				wg.Done()
				return
			}
		}
	}()
}
