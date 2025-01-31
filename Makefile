.PHONY: build clean producer consumer deps-up test gen-pb-go

build:
	go build -o bin/consumer consumer/*.go
	go build -o bin/producer producer/*.go

deps-up:
	make -f ../docker360/kubernetes/Makefile system-start
	make -f ../docker360/kubernetes/Makefile kafka-environment-start

docker-deps-up:
	docker-compose up -d

docker-deps-down:
	docker-compose down

gen-pb-go:
	protoc *.proto --go_out=.

producer:
	go run producer/*.go

consumer:
	go run consumer/*.go

producer-confluent:
	ENV=confluent go run producer/*.go

consumer-confluent:
	ENV=confluent go run consumer/*.go

producer-docker:
	ENV=docker go run producer/*.go

consumer-docker:
	ENV=docker go run consumer/*.go

test:
	curl --location --request POST 'localhost:3000/api/v1/account-delete' --header 'Content-Type: application/json' --data-raw '{ "userId":"2016fe16-4e40-4b3c-87a2-3675ff1f8d97", "reason":"deleted" }'
