# kafka-with-go

Kafka Producer and Consumer for proto messages. 

The Producer exposes an endpoint `localhost:3000/api/v1/account-delete` to publish to `life360_account_deleted` topic.

Reference: https://medium.com/swlh/apache-kafka-with-golang-227f9f2eb818

## Requirements

Kafka running on `localhost:32092`

## How to run with minikube

1. `make build` to build binaries
2. `make deps-up` to run kafka
3. `make producer` to run the producer (starts HTTP server)
4. `make consumer` to run the consumer
5. `make test` to send a message to the producer

## How to run with docker-compose

1. `make deps-up` to run kafka
2. `make producer-docker` to run the producer (starts HTTP server)
3. `make consumer-docker` to run the consumer
4. `make test` to send a message to the producer

## How to run with confluent

Update `config/config-confluent.properties` with real values. Go to https://confluent.cloud/settings/api-keys/create to create API keys:
- `bootstrap.servers` is the URL of your Kafka cluster
- `sasl.username` is your API key
- `sasl.password` is your API secret

1. `make producer-confluent` to run the producer (starts HTTP server)
2. `make consumer-confluent` to run the consumer
3. `make test` to send a message to the producer


## Protobuf setup

```
# install protoc
brew install protobuf

# install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# install go protobuf library
go get google.golang.org/protobuf

# complie protobuf (generates life360_account_deleted_v1.pb.go)
protoc *.proto --go_out=.
```

## Links:
- https://github.com/confluentinc/demo-scene/blob/master/getting-started-with-ccloud-golang/ClientApp.go
