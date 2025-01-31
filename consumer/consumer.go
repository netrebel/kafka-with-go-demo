package main

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/life360/kafka-with-go-demo/config"
	"github.com/life360/kafka-with-go-demo/protos"
	"google.golang.org/protobuf/proto"
)

func connectConsumer(brokersURL string) (*kafka.Consumer, error) {
	p, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               brokersURL,
		"group.id":                        "foo",
		"go.application.rebalance.enable": true})

	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created consumer on %v\n", brokersURL)
	return p, nil
}

func main() {
	config.LoadEnv()
	topic := os.Getenv("TOPIC")
	consumer, err := connectConsumer(os.Getenv("bootstrap.servers"))

	if err != nil {
		panic(err)
	}

	// Calling ConsumePartition. It will open one connection per broker
	// and share it for all partitions that live on it.
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Subscibed to topic: %v\n", topic)

	run := true
	for run {
		ev := consumer.Poll(10)
		switch e := ev.(type) {
		case *kafka.Message:
			msg := &protos.Life360AccountDeleted{}
			err = proto.Unmarshal(e.Value, msg)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Message on %v: %v\n", e.TopicPartition, msg)
		case kafka.PartitionEOF:
			fmt.Printf("Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		default:
			// fmt.Printf("Ignored %v\n", e)
			// run = false
		}
	}

	consumer.Close()
}
