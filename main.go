package main

import (
	ogiconsumer "github.com/gojekfarm/ogi/consumer"
	logger "github.com/gojekfarm/ogi/logger"
)

func main() {
	logger.SetupLogger()
	ogiconsumer.Consume()
}
