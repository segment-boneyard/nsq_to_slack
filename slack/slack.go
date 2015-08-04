package slack

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/bitly/go-nsq"
	"github.com/segmentio/go-log"
)

// Slack structure.
type Slack struct {
	url string
}

// New initializes Slack.
func New(url string) *Slack {
	return &Slack{url}
}

// HandleMessage decodes the message into `Message`
// validates it and sends it to slack, when validation
// error occurs the message is discarded.
func (s *Slack) HandleMessage(msg *nsq.Message) error {
	r := bytes.NewReader(msg.Body)
	resp, err := http.Post(s.url, "application/json", r)
	if err != nil {
		return nil
	}

	switch resp.StatusCode/100 | 0 {
	case 5:
		log.Warning("server error: %s", resp.Status)
		return fmt.Errorf("server error: %s", resp.Status)
	case 4:
		buf, _ := httputil.DumpResponse(resp, true)
		log.Error("%s:\n%s", resp.StatusCode, buf)
		return nil
	}

	return nil
}
