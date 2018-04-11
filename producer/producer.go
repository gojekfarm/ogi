package ogiproducer

import (
	"fmt"

	"github.com/abhishekkr/gol/golenv"
	"github.com/gojekfarm/ogi/instrumentation"
	"github.com/gojekfarm/ogi/logger"
)

type Producer interface {
	NewProducer()
	Close()
	GetMetadata()
	GetPartitionNumber(topic string, messageKey string) (partitionNumber int32)
	ProduceMessage(topic string, message []byte, partitionNumber int32)
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

func Produce(producer Producer, topic string, message []byte, message_key string) {
	txn := instrumentation.StartTransaction("produce_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	partitionNumber := producer.GetPartitionNumber(topic, message_key)
	producer.ProduceMessage(topic, message, partitionNumber)
	logger.Infof("topic '%s' message-key '%s' is assigned to '%d' partition", topic, message_key, partitionNumber)
}
