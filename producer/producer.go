package main

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/gofiber/fiber/v2"
)

// Comment struct
type Comment struct {
	Text string `form:"text" json:"text"`
}

func createComment(c *fiber.Ctx) error {
	// Instantiate new Message struct
	cmt := new(Comment)
	if err := c.BodyParser(cmt); err != nil {
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}
	// convert body into bytes and send it to kafka
	cmtInBytes, err := json.Marshal(cmt)
	if err != nil {
		fmt.Printf("Error marshalling comment: %s", err)
		return err
	}

	err = PushCommentToQueue("comments", cmtInBytes)
	if err != nil {
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Error creating comment",
		})
		return err
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Comment pushed successfully",
		"comment": cmt,
	})
}

// PushCommentToQueue pushes commen to the topic
func PushCommentToQueue(topic string, message []byte) error {
	brokersURL := []string{"localhost:29092"}
	fmt.Printf("Connecting to Kafka on: %v\n", brokersURL)
	producer, err := ConnectProducer(brokersURL)
	if err != nil {
		fmt.Printf("Error connecting to producer: %v", err)
		return err
	}
	defer producer.Close()
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	fmt.Printf("Message received: %+v\n", msg)
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Printf("Error sending message to producer: %v", err)
		return err
	}
	fmt.Printf("Message is stored in: topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

// ConnectProducer connects to Kafka
func ConnectProducer(brokersURL []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	// NewSyncProducer creates a new SyncProducer using the given broker addresses and configuration.
	conn, err := sarama.NewSyncProducer(brokersURL, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {

	app := fiber.New()
	api := app.Group("/api/v1") // /api

	api.Post("/comments", createComment)

	fmt.Println("Start sending messages to localhost:3000/api/v1/comments")

	app.Listen(":3000")

}
