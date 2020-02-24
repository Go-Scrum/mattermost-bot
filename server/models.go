package main

import "time"

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
	WorkspaceID string       `db:"workspace_id" json:"workspace_id"`
	ChannelID   string       `db:"channel_id" json:"channel_id"`
	UserID      string       `db:"user_id" json:"user_id"`
	QuestionID  string       `db:"question_id" json:"question_id"`
	Comment     string       `db:"comment" json:"comment"`
	Status      AnswerStatus `db:"status" json:"status"`
	MessageTS   string       `db:"message_ts" json:"message_ts"`
	Question    Question
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

type QuestionType int

const (
	Text       QuestionType = 0
	Numeric    QuestionType = 1
	PreDefined QuestionType = 2
)

type Question struct {
	Model
	Title     string
	Type      QuestionType
	ProjectId string
}

type WorkspaceType int

const (
	Mattermost WorkspaceType = 0
)

type Workspace struct {
	Model
	BotUserID     string        `db:"bot_user_id" json:"bot_user_id"`
	Language      string        `db:"language" json:"language"`
	WorkspaceName string        `db:"workspace_name" json:"workspace_name" `
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
