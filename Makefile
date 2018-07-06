.PHONY: all
all: build test


setup:
	mkdir -p ${GOPATH}/bin && \
	if which dep &> /dev/null ; then go get -u github.com/golang/dep/cmd/dep ; fi


build-deps: setup
	dep ensure

compile: build-deps
	set -e ; \
	mkdir -p out

	set -e ; \
	export GOOS=linux ; export GOARCH=amd64 ; \
		go build -o out/ogi main.go ; \
		echo "[info] compiled ogi main"

	set -e ; \
	export GOOS=linux ; export GOARCH=amd64 ; \
		cd plugin-examples/transformer/message_logs ; \
		go build -o "../../../out/transformer-message-log.so" -buildmode=plugin . ; \
		cd - ; \
		echo "[info] compiled transformer.message_logs plugin"

	set -e ; \
	export GOOS=linux ; export GOARCH=amd64 ; \
		cd plugin-examples/transformer/os_path_exists ; \
		go build -o "../../../out/transformer-os-path-exists.so" -buildmode=plugin . ; \
		cd - ; \
		echo "[info] compiled transformer.os_path_exists plugin"

	set -e ; \
	export GOOS=linux ; export GOARCH=amd64 ; \
		cd plugin-examples/producer/echo ; \
		go build -o "../../../out/producer-echo.so" -buildmode=plugin . ; \
		cd - ; \
		echo "[info] compiled producer.echo plugin"

	set -e ; \
	export GOOS=linux ; export GOARCH=amd64 ; \
		cd plugin-examples/producer/filedump ; \
		go build -o "../../../out/producer-filedump.so" -buildmode=plugin .; \
		cd - ; \
		echo "[info] compiled producer.filedump plugin"

	set -e ; \
	export GOOS=linux ; export GOARCH=amd64 ; \
		cd plugin-examples/consumer/gcp_stackdriver_logs ; \
		go build -o "../../../out/consumer-gcp-stackdriver-logs.so" -buildmode=plugin . ; \
		cd - ; \
		echo "[info] compiled consumer.gcp_stackdriver_logs plugin"

	set -e ; \
	export GOOS=linux ; export GOARCH=amd64 ; \
		cd plugin-examples/consumer/file_line_by_line ; \
		 go build -o "../../../out/consumer-file-line-by-line.so" -buildmode=plugin . ; \
		 cd - ; \
		echo "[info] compiled consumer.file_line_by_line plugin"

	echo "done."

build: build-deps compile

build-test-plugins: build-deps
	set -e ; \
	export THIS_DIR=$(pwd) ; \
		cd tests/consumer ; \
		go build -o "../consumer.so" -buildmode=plugin . ; cd -
	set -e ; \
	export THIS_DIR=$(pwd) ; \
		cd tests/consumer-bad ; \
		go build -o "../consumer-bad.so" -buildmode=plugin . ; cd -
	set -e ; \
	export THIS_DIR=$(pwd) ; \
		cd tests/transformer ; \
		go build -o "../transformer.so" -buildmode=plugin . ; cd -
	set -e ; \
	export THIS_DIR=$(pwd) ; \
		cd tests/transformer-bad ; \
		go build -o "../transformer-bad.so" -buildmode=plugin . ; cd -
	set -e ; \
	export THIS_DIR=$(pwd) ; \
		cd tests/producer ; \
		go build -o "../producer.so" -buildmode=plugin . ; cd -
	set -e ; \
	export THIS_DIR=$(pwd) ; \
		cd tests/producer-bad ; \
		go build -o "../producer-bad.so" -buildmode=plugin . ; cd -

test: build-test-plugins
	go test -gcflags=-l github.com/gojektech/ogi/consumer github.com/gojektech/ogi/transformer github.com/gojektech/ogi/producer
