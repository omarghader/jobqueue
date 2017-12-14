package publisher

import (
	"github.com/Sirupsen/logrus"
	"github.com/omarghader/jobqueue/protocol"
	"github.com/smallnest/goreq"
)

func (s *publisherService) PublishWebhook(job *protocol.Event) error {

	for _, url := range s.webhookURLs {
		_, _, err := goreq.New().Post(url).SendStruct(job).End()
		if err != nil {
			logrus.Errorln("Cannot Send Webhook notification to :", url)
		}
	}
	return nil
}
