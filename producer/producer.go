package ogiproducer

import (
	"fmt"

	"github.com/abhishekkr/gol/golenv"
	"github.com/gojekfarm/kafka-ogi/instrumentation"
	"github.com/gojekfarm/kafka-ogi/logger"
)

type Producer interface {
	NewProducer()
	Close()
	GetMetadata()
	GetPartitionNumber(topic string, messageKey string) (partitionNumber int32)
	ProduceMessage(topic string, message []byte, partitionNumber int32)
}

var (
	BootstrapServers = golenv.OverrideIfEnv("PRODUCER_BOOTSTRAP_SERVERS", "")
	PartitionCounts  map[string]int
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
	k := &Kafka{}
	k.NewProducer()
	return k
}

func Produce(producer Producer, topic string, message []byte, message_key string) {
	txn := instrumentation.StartTransaction("produce_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	partitionNumber := producer.GetPartitionNumber(topic, message_key)
	producer.ProduceMessage(topic, message, partitionNumber)
	logger.Infof("topic '%s' message-key '%s' is assigned to '%d' partition", topic, message_key, partitionNumber)
}
