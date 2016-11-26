configFile=$(shell pwd)/config.yaml

help:
	@echo "Available commands:"
	@echo " - make build - compile project and dependencies"
	@echo " - make test 	- run all tests"
	@echo " - make clean - remove bin files"

build:
	go get github.com/julienschmidt/httprouter
	go get github.com/nlopes/slack
	go get gopkg.in/olivere/elastic.v3
	go get gopkg.in/yaml.v2
	go get github.com/satori/go.uuid
	go get github.com/Sirupsen/logrus
	go get github.com/kennygrant/sanitize

	go build -o slackbot-links main.go

test:
	go test ./src/links/ -config-file=$(configFile)
	go test ./src/mercury/ -config-file=$(configFile)

docapi:
	[ ! -f snowboard ] && curl -L https://github.com/subosito/snowboard/releases/download/v0.4.3/snowboard-v0.4.3.linux-amd64.tar.gz | tar -xz || true
	./snowboard -i API.apib -o docs/index.html

clean:
	rm -f slackbot-links
