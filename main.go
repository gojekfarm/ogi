package main

import (
	ogiconsumer "github.com/gojekfarm/kafka-ogi/consumer"
	logger "github.com/gojekfarm/kafka-ogi/logger"
)

func main() {
	logger.SetupLogger()
	var k ogiconsumer.Kafka
	ogiconsumer.Consume(&k)
}
