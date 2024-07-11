package handler

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2/log"
	"github.com/nsqio/go-nsq"
)

type Notification struct {
	Title    string  `json:"title"`
	Body     string  `json:"body"`
	ImageUrl *string `json:"image_url,omitempty"`
}

func (n *Notification) Byte() []byte {
	b, _ := json.Marshal(n)
	return b
}

func (n Notification) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	if err := json.Unmarshal(m.Body, &n); err != nil {
		log.Errorf("job id %s failed to marshal json: %v", m.ID, err)
		return err
	}

	log.Infof("received notification-%s with body %v", m.ID, n)

	return nil
}
