package main

import (
	"context"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
	"worker-mesh/config"
	"worker-mesh/internal/handler"
	"worker-mesh/pkg/messaging"
	"worker-mesh/pkg/router"

	"github.com/gofiber/fiber/v2/log"
	"github.com/tjarratt/babble"
)

func init() {
	config.InitConfig()

	if config.ReadConfig().Env == "production" {
		log.SetLevel(log.LevelError)
	}
}

func main() {
	producer, stop := messaging.NewProducer()
	defer stop()
	go emulateRequests(producer)

	httpHandler := handler.NewHttpHandler(producer)

	app := router.NewFiber(httpHandler)

	go func() {
		if err := app.Listen(":" + config.ReadConfig().Port); err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdown, release := context.WithTimeout(context.Background(), 5*time.Second)
	defer release()

	if err := app.ShutdownWithContext(shutdown); err != nil {
		log.Fatal(err)
	}
}

func emulateRequests(producer messaging.Producer) {
	max, min := 3, 1
	maxDur, minDur := 100, 1

	babbler := babble.NewBabbler()
	babbler.Separator = " "

	url := "https://picsum.photos/300/200"

	if producer.Ping() != nil {
		log.Fatal("failed to ping nsq")
	}

	for {
		body := handler.Notification{
			Title:    babbler.Babble(),
			Body:     babbler.Babble(),
			ImageUrl: &url,
		}

		if rand.IntN(3) == 1 {
			dur := rand.IntN(maxDur+1-minDur) + minDur
			delay := time.Duration(dur) * time.Second

			if err := producer.DeferredPublish(config.ReadConfig().Nsq.Topic, delay, body.Byte()); err != nil {
				log.Errorf("failed to publish emulated message: %v", err)
			}

			log.Infof("sent message: %v with delay: %d second", body, dur)

		} else {
			if err := producer.Publish(config.ReadConfig().Nsq.Topic, body.Byte()); err != nil {
				log.Errorf("failed to publish emulated message: %v", err)
			}

			log.Infof("sent message: %v", body)
		}

		time.Sleep(time.Duration(rand.IntN(max+1-min)+min) * time.Second)
	}
}
