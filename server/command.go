package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "goscrum",
		DisplayName:      "goscrum",
		Description:      "GoScrum for your standup", // TODO -- need to fix the name
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: standup, help",
		AutoCompleteHint: "[command]",
	}
}

func (p *Plugin) postCommandResponse(args *model.CommandArgs, text string) {
	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: args.ChannelId,
		Message:   text,
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)
}

// ExecuteCommand command help text
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	command := split[0]
	//parameters := []string{}
	//action := ""
	//if len(split) > 1 {
	//	action = split[1]
	//}
	//if len(split) > 2 {
	//	parameters = split[2:]
	//}

	if command != "/goscrum" {
		return &model.CommandResponse{}, nil
	}

	text := "###### Mattermost goscrum Plugin - Slash Command Help\n" + strings.Replace(`* |/welcomebot preview [team-name] [user-name]| - preview the welcome message for the given team name. The current user's username will be used to render the template.
* |/welcomebot list| - list the teams for which welcome messages were defined`, "|", "`", -1)
	p.postCommandResponse(args, text)

	//switch action {
	//case "standup":
	//	var str strings.Builder
	//	str.WriteString("Welcome to goscrum")
	//	p.postCommandResponse(args, str.String())
	//case "help":
	//	fallthrough
	//case "":
	//	text := "###### Mattermost welcomebot Plugin - Slash Command Help\n" + strings.Replace(COMMAND_HELP, "|", "`", -1)
	//	p.postCommandResponse(args, text)
	//	return &model.CommandResponse{}, nil
	//}
	return &model.CommandResponse{}, nil
}
