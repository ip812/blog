package notifier

import (
	"github.com/ip812/blog/logger"
	"github.com/slack-go/slack"
)

type Slack struct {
	api *slack.Client
	log logger.Logger
}

func NewSlack(token string, log logger.Logger) *Slack {
	return &Slack{
		api: slack.New(token),
		log: log,
	}
}

func (s *Slack) SendMsg(
	channelID string,
	text string,
) error {
	_, _, err := s.api.PostMessage(
		channelID,
		slack.MsgOptionText(text, false),
	)
	if err != nil {
		s.log.Error("failed to send message %v to Slack channel: %s", err, channelID)
		return err
	}
	s.log.Info("Message sent successfully to Slack channel: %s", channelID)

	return nil
}
