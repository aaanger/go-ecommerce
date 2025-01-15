package main

import (
	"ecommerce/pkg/broker"
	"github.com/sirupsen/logrus"
)

func main() {
	topics := []string{"orders.create", "orders.update", "orders.cancel"}

	consumer, err := broker.NewConsumer("localhost:9092", topics)
	if err != nil {
		logrus.Fatalf("Error initializing kafka conumser: %s", err)
	}

	go consumer.Run()

	err = consumer.Stop()
	if err != nil {
		logrus.Errorf("Error stopping kafka consumer: %s", err)
	}
}
