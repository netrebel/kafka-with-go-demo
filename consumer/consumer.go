package main

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/netrebel/kafka-with-go/protos"
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
	topic := "life360_account_deleted"
	consumer, err := connectConsumer("192.168.64.26:32092")

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
	for run == true {
		ev := consumer.Poll(10)
		switch e := ev.(type) {
		case *kafka.Message:
			msg := &protos.Life360AccountDeleted{}
			err = proto.Unmarshal(e.Value, msg)
			if err != nil {
				panic(err)
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
