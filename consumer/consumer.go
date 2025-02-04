package main

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/life360/kafka-with-go-demo/config"
	"github.com/life360/kafka-with-go-demo/protos"
	"google.golang.org/protobuf/proto"
)

func connectConsumer() (*kafka.Consumer, error) {
	fmt.Printf("Connecting to Kafka on: %v\n", os.Getenv("bootstrap.servers"))

	config := kafka.ConfigMap{
		"bootstrap.servers":               os.Getenv("bootstrap.servers"),
		"group.id":                        "go-group-1",
		"auto.offset.reset":               "latest",
		"session.timeout.ms":              45000, // Best practice for higher availability in librdkafka clients prior to 1.7
		"go.application.rebalance.enable": true,
		"client.id":                       "ccloud-go-client-local",
	}

	if os.Getenv("sasl.username") != "" && os.Getenv("sasl.password") != "" {
		config["security.protocol"] = "SASL_SSL"
		config["sasl.mechanisms"] = "PLAIN"
		config["sasl.username"] = os.Getenv("sasl.username")
		config["sasl.password"] = os.Getenv("sasl.password")
	}

	p, err := kafka.NewConsumer(&config)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	return p, nil
}

func main() {
	config.LoadEnv()

	topic := os.Getenv("env.topic")
	if topic == "" {
		panic("env.topic not set")
	}

	consumer, err := connectConsumer()

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
