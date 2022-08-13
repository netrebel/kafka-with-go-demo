.PHONY: build clean producer consumer

build:
	cd consumer; go build -o ../bin/consumer .
	cd producer; go build -o ../bin/producer .

producer:
	bin/producer

consumer:
	bin/consumer