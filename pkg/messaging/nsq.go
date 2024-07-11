package messaging

import (
	"fmt"
	"log"
	"os"
	"time"
	"worker-mesh/config"

	"github.com/nsqio/go-nsq"
)

type Producer interface {
	Publish(topic string, body []byte) error
	DeferredPublish(topic string, delay time.Duration, body []byte) error

	Ping() error
}

func NewProducer() (producer *nsq.Producer, stop func()) {
	cfg := nsq.NewConfig()
	cfg.Hostname = fmt.Sprintf("producer-%v", os.Getpid())

	producer, err := nsq.NewProducer(config.ReadConfig().Nsq.Address(), cfg)
	if err != nil {
		log.Fatalf("failed to init nsq producer: %v", err)
	}

	return producer, producer.Stop
}

func NewConsumer(topic string, channel string, handlers ...nsq.Handler) (stop func()) {
	cfg := nsq.NewConfig()
	cfg.Hostname = fmt.Sprintf("consumer-%v", os.Getpid())

	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		log.Fatalf("failed to init nsq consumer: %v", err)
	}

	for _, h := range handlers {
		consumer.AddHandler(h)
	}

	if err := consumer.ConnectToNSQLookupd(config.ReadConfig().Nsq.Lookupd.Address()); err != nil {
		log.Fatalf("failed to connect to nsqlookupd: %v", err)
	}

	return consumer.Stop
}
