package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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

	topic := os.Getenv("env.topic")
	if topic == "" {
		panic("env.topic not set")
	}

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
	p, err := ConnectProducer()
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
func ConnectProducer() (*kafka.Producer, error) {
	fmt.Printf("Connecting to Kafka on: %v\n", os.Getenv("bootstrap.servers"))

	config := kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("bootstrap.servers"),
		"client.id":         "ccloud-go-client-2f0d4f57-0582-4c6d-8442-24513c4f715a",
		"acks":              "all",
	}
	if os.Getenv("sasl.username") != "" && os.Getenv("sasl.password") != "" {
		config["security.protocol"] = "SASL_SSL"
		config["sasl.mechanisms"] = "PLAIN"
		config["sasl.username"] = os.Getenv("sasl.username")
		config["sasl.password"] = os.Getenv("sasl.password")
	}
	p, err := kafka.NewProducer(&config)

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
