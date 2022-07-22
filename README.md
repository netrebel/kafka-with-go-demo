# kafka-with-go

https://medium.com/swlh/apache-kafka-with-golang-227f9f2eb818

## How to run

* `docker-compose up -d` to run kafka
* `go run producer/*.go` to run the producer (starts HTTP server)
* `protoc *.proto --go_out=:$GOPATH/src` to compile protobuf file
* `curl --location --request POST 'localhost:3000/api/v1/account-delete' --header 'Content-Type: application/json' --data-raw '{ "userId":"2016fe16-4e40-4b3c-87a2-3675ff1f8d97", "reason":"deleted" }'` to send a message to the producer
* `go run consumer/consumer.go` to run the consumer


## Protobuf setup

```
# install protoc
brew install protobuf

# install go protobuf 
go get github.com/golang/protobuf
go get github.com/golang/protobuf/protoc-gen-go

# complie protobuf
protoc *.proto --go_out=producer
```