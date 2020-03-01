package main

import (
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
)

type Model struct {
	ID        string    `json:"id" valid:"uuidv4, optional"`
	CreatedAt time.Time `valid:"-" json:"created_at"`
	UpdatedAt time.Time `valid:"-" json:"-"`
}

type AnswerStatus string

const (
	Asked     AnswerStatus = "asked"
	Answered  AnswerStatus = "answered"
	Cancelled AnswerStatus = "cancelled"
)

// Standup model used for serialization/deserialization stored standups
type Answer struct {
	Model
	ParticipantID string `json:"participant_id"`
	QuestionID    string `json:"question_id"`
	Comment       string `json:"comment"`
	BotPostId     string `json:"bot_post_id"`
	Question      Question
	Participant   Participant
}

type Participant struct {
	Model
	WorkspaceID string `db:"workspace_id" json:"workspace_id"`
	UserID      string `db:"user_id" json:"user_id"`
	Role        string `db:"role" json:"role"`
	RealName    string `db:"real_name" json:"real_name"`
	FirstName   string `db:"first_name" json:"first_name"`
	LastName    string `db:"last_name" json:"last_name"`
	NickName    string `db:"nick_name" json:"nick_name"`
	Email       string `db:"email" json:"channel_name"`
	Projects    []*Project
}

type Project struct {
	Model
	WorkspaceID      string `json:"workspace_id"`
	ChannelName      string `json:"channel_name"`
	Name             string `json:"name"`
	ChannelID        string `json:"channel_id"`
	Deadline         string `json:"deadline"`
	TZ               string `json:"tz"`
	OnbordingMessage string `json:"onbording_message,omitempty"`
	SubmissionDays   string `json:"submission_days,omitempty"`
	Participants     []*Participant
	ReportingChannel string `json:"reporting_channel"`
	ReportingTime    string `json:"reporting_time"`
	IsActive         bool   `json:"is_active"`
	Questions        []*Question
}

type QuestionType string

const (
	Text       QuestionType = "Text"
	Numeric    QuestionType = "Numeric"
	PreDefined QuestionType = "PreDefined"
)

type Question struct {
	Model
	Title     string
	Type      QuestionType
	ProjectId string
}

type WorkspaceType string

const (
	Mattermost WorkspaceType = "Mattermost"
)

type Workspace struct {
	Model
	BotUserID     string        `json:"bot_user_id"`
	Language      string        ` json:"language"`
	WorkspaceName string        `json:"workspace_name" `
	URL           string        `json:"url"`
	WorkspaceType WorkspaceType `json:"workspace_type"`
	ClientID      string        `json:"client_id"`
	ClientSecret  string        `json:"client_secret"`
	AccessToken   string        `json:"access_token"`
	RefreshToken  string        `json:"refresh_token"`
	Expiry        *time.Time    `json:"expiry,omitempty"`
	PersonalToken string        `json:"personal_token,omitempty"`
	Projects      []*Project
}

type MessageType string

const (
	QuestionMessage   MessageType = "QuestionMessage"
	AnswerMessage     MessageType = "AnswerMessage"
	StandupMessage    MessageType = "StandupMessage"
	ReportMessage     MessageType = "ReportMessage"
	GreetingMessage   MessageType = "GreetingMessage"
	OnBoardingMessage MessageType = "OnBoardingMessage"
)

type Message struct {
	Model
	Attachments   []*model.SlackAttachment
	Content       string      `json:"content"`
	UserId        string      `json:"user_id"`
	ChannelId     string      `json:"channel_id"`
	MessageType   MessageType `json:"message_type"`
	ParticipantID string      `json:"participant_id"`
	Question      Question
	Participant   Participant
}

type UserActivityType string

const (
	UserQuestionActivity   UserActivityType = "UserQuestionActivity"
	UserAnswerActivity     UserActivityType = "UserAnswerActivity"
	UserStandupActivity    UserActivityType = "UserStandupActivity"
	UserReportActivity     UserActivityType = "UserReportActivity"
	UserGreetingActivity   UserActivityType = "UserGreetingActivity"
	UserOnBoardingActivity UserActivityType = "UserOnBoardingActivity"
)

type UserActivity struct {
	UserId        string           `json:"user_id"`
	ChannelID     string           `json:"channel_id"`
	ProjectID     string           `json:"project_id"`
	ParticipantID string           `json:"participant_id"`
	QuestionID    string           `json:"question_id"`
	BotPostId     string           `json:"bot_post_id"`
	ActivityType  UserActivityType `json:"activity_type"`
	Question      Question
	Participant   Participant
	Project       Project
}
