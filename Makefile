.PHONY: all
all: build test


setup:
	if which dep &> /dev/null ; then go get -u github.com/golang/dep/cmd/dep ; fi


build-deps: setup
	dep ensure

compile: build-deps
	mkdir -p out

	GOOS=linux GOARCH=amd64 go build -o out/ogi main.go

	cd plugin-examples/transformer/message_logs ; \
		go build -o "../../../out/transformer-message-log.so" -buildmode=plugin . ; \
		cd -
	cd plugin-examples/producer/echo; \
		go build -o "../../../out/producer-echo.so" -buildmode=plugin . ; \
		cd -
	cd plugin-examples/producer/filedump; \
		go build -o "../../../out/producer-filedump.so" -buildmode=plugin . ; \
		cd -
	cd plugin-examples/consumer/gcp_stackdriver_logs; \
		go build -o "../../../out/consumer-gcp-stackdriver-logs.so" -buildmode=plugin . ; \
		cd -

build: build-deps compile

build-test-plugins:
	export THIS_DIR=$(pwd)
	cd tests/consumer ; \
		go build -o "../consumer.so" -buildmode=plugin . ; cd -
	cd tests/consumer-bad ; \
		go build -o "../consumer-bad.so" -buildmode=plugin . ; cd -
	cd tests/transformer ; \
		go build -o "../transformer.so" -buildmode=plugin . ; cd -
	cd tests/transformer-bad ; \
		go build -o "../transformer-bad.so" -buildmode=plugin . ; cd -
	cd tests/producer ; \
		go build -o "../producer.so" -buildmode=plugin . ; cd -
	cd tests/producer-bad ; \
		go build -o "../producer-bad.so" -buildmode=plugin . ; cd -

test: build-test-plugins
	go test -gcflags=-l github.com/gojekfarm/ogi/consumer github.com/gojekfarm/ogi/transformer github.com/gojekfarm/ogi/producer
