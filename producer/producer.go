package ogiproducer

import (
	"fmt"

	"github.com/abhishekkr/gol/golenv"
	"github.com/gojekfarm/ogi/instrumentation"
	logger "github.com/gojekfarm/ogi/logger"
)

type Producer interface {
	NewProducer()
	Close()
	Produce(string, []byte, string)
}

type NewProducerFunc func() Producer

var (
	BootstrapServers = golenv.OverrideIfEnv("PRODUCER_BOOTSTRAP_SERVERS", "")
	ProducerType     = golenv.OverrideIfEnv("PRODUCER_TYPE", "confluent-kafka")

	producerMap = map[string]NewProducerFunc{
		"confluent-kafka": NewConfluentKafka,
	}
)

func validateConfig() {
	var missingVariables string
	if BootstrapServers == "" {
		missingVariables = fmt.Sprintf("%s PRODUCER_BOOTSTRAP_SERVERS", missingVariables)
	}

	if missingVariables != "" {
		logger.Fatalf("Missing Env Config:%s", missingVariables)
	}
}

func NewProducer() Producer {
	return producerMap[ProducerType]()
}

func Produce(topic string, message []byte, messageKey string) {
	txn := instrumentation.StartTransaction("produce_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	producer := NewProducer()
	defer producer.Close()

	producer.Produce(topic, message, messageKey)
	logger.Infof("topic '%s' message-key '%s'", topic, messageKey)
}
