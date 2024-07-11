package main

import (
	"os"
	"os/signal"
	"syscall"
	"worker-mesh/config"
	"worker-mesh/internal/handler"
	"worker-mesh/pkg/messaging"
)

func init() {
	config.InitConfig()
}

func main() {
	stop := messaging.NewConsumer(
		config.ReadConfig().Nsq.Topic,
		config.ReadConfig().Nsq.Channel,
		handler.Notification{},
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	stop()
}
