# kafka-with-go

https://medium.com/swlh/apache-kafka-with-golang-227f9f2eb818

## How to run

* `docker-compose up -d` to run kafka
* `go run producer/producer.go` to run the producer (starts HTTP server)
* `curl --location --request POST 'localhost:3000/api/v1/comments' --header 'Content-Type: application/json' --data-raw '{ "text":"nice boy" }'` to send a message to the producer