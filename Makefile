configFile=$(shell pwd)/config.yaml

help:
	@echo "Available commands:"
	@echo " - make build  - compile project and dependencies"
	@echo " - make test   - run all tests"
	@echo " - make docapi - generate html documentation from API.apib file"
	@echo " - make clean  - remove bin files"

build:
	go get -d -v ./...
	go build -o slackbot-links main.go

test:
	go test ./src/links/ -config-file=$(configFile) ${ARGS}
	go test ./src/mercury/ -config-file=$(configFile) ${ARGS}

docapi:
	[ ! -f snowboard ] && curl -L https://github.com/subosito/snowboard/releases/download/v0.4.3/snowboard-v0.4.3.linux-amd64.tar.gz | tar -xz || true
	./snowboard -i API.apib -o docs/index.html

clean:
	rm -f slackbot-links
