package slackrus

import (
	"github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

const (
	VERSION = "1.0.0"
)

const (
	SLACK_COLOR_DANGER  string = "danger"
	SLACK_COLOR_WARNING string = "warning"
	SLACK_COLOR_GOOD    string = "good"
)

type SlackHook struct {
	client     *slack.Client
	token      string
	Username   string
	AuthorName string
	Channel    string
	IconURL    string
	IconEmoji  string
}

func NewSlackHook(token, username, authorName, channel, iconURL, iconEmoji string) *SlackHook {
	hook := &SlackHook{
		token:      token,
		Username:   username,
		AuthorName: authorName,
		Channel:    channel,
		IconURL:    iconURL,
		IconEmoji:  iconEmoji,
	}

	return hook
}

func (hook *SlackHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		//		logrus.WarnLevel,
		//		logrus.InfoLevel,
		//		logrus.DebugLevel,
	}
}

func (hook *SlackHook) Fire(sourceEntry *logrus.Entry) error {
	hook.client = slack.New(hook.token)

	params := slack.PostMessageParameters{
		Username:  hook.Username,
		IconURL:   hook.IconURL,
		IconEmoji: hook.IconEmoji,
	}

	var messageFields []slack.AttachmentField

	for key, value := range sourceEntry.Data {
		message := slack.AttachmentField{
			Title: key,
			Value: value.(string),
			Short: true,
		}

		messageFields = append(messageFields, message)
	}

	attachment := slack.Attachment{
		Color:      getColor(sourceEntry.Level),
		AuthorName: hook.AuthorName,
		Fields:     messageFields,
		Text:       sourceEntry.Message,
	}

	params.Attachments = []slack.Attachment{attachment}
	_, _, err := hook.client.PostMessage(hook.Channel, "", params)

	return err
}

func getColor(level logrus.Level) string {
	switch level {
	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
		return SLACK_COLOR_DANGER
	case logrus.WarnLevel:
		return SLACK_COLOR_WARNING
	case logrus.InfoLevel:
		return SLACK_COLOR_GOOD
	default:
		return ""
	}
}

