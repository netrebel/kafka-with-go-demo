.PHONY: build clean producer consumer

build:
	go build -o bin/consumer consumer/*.go
	go build -o bin/producer producer/*.go

producer:
	bin/producer

consumer:
	bin/consumer