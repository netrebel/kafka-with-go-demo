package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
	"github.com/life360/kafka-with-go-demo/config"
	"github.com/life360/kafka-with-go-demo/protos"
	proto "google.golang.org/protobuf/proto"
)

func createMessage(w http.ResponseWriter, r *http.Request) {
	// read body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: fail to read request body: %s", err)
		response(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("INFO: JSON Request: %v", string(data))

	// unmarshal and create docMsg
	msg := &protos.Life360AccountDeleted{}
	err = json.Unmarshal(data, msg)
	if err != nil {
		log.Printf("ERROR: fail unmarshl: %s", err)
		response(w, "Invalid request json", 400)
	}

	log.Printf("INFO: Life360AccountDeleted: %v", msg)

	protoMsg, err := proto.Marshal(msg)
	if err != nil {
		log.Fatalln("Failed to proto encode doc:", err)
	}

	topic := os.Getenv("TOPIC")
	err = PushMessageToTopic(topic, protoMsg)
	if err != nil {
		http.Error(w, "Error pushing to topic", http.StatusInternalServerError)
	} else {
		fmt.Printf("Published message to topic %v\n", topic)
	}

}

func response(w http.ResponseWriter, resp string, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, resp)
}

// PushMessageToTopic pushes commen to the topic
func PushMessageToTopic(topic string, message []byte) error {
	brokersURL := os.Getenv("bootstrap.servers")
	fmt.Printf("Connecting to Kafka on: %v\n", brokersURL)
	p, err := ConnectProducer(brokersURL)
	if err != nil {
		fmt.Printf("Error connecting producer: %v", err)
		return err
	}
	defer p.Close()

	delivery_chan := make(chan kafka.Event, 10000)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message},
		delivery_chan,
	)
	if err != nil {
		panic("Could not produce message")
	}
	p.Flush(1000)
	return nil
}

// ConnectProducer connects to Kafka
func ConnectProducer(brokersURL string) (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokersURL,
		"client.id":         "local",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	return p, nil
}

func main() {

	// router
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/account-delete", createMessage).Methods("POST")
	config.LoadEnv()

	log.Printf("Start sending messages to localhost:3000/api/v1/account-delete")

	// start server
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Printf("ERROR: fail init http server, %s", err)
		os.Exit(1)
	}

}
