package ogiproducer

import (
	"fmt"

	"github.com/abhishekkr/gol/golenv"
	"github.com/gojekfarm/ogi/instrumentation"
	logger "github.com/gojekfarm/ogi/logger"
)

type Producer interface {
	Produce(string, []byte, string)
	Close()
}

type NewProducerFunc func() Producer

var (
	BootstrapServers = golenv.OverrideIfEnv("PRODUCER_BOOTSTRAP_SERVERS", "")
	ProducerType     = golenv.OverrideIfEnv("PRODUCER_TYPE", "confluent-kafka")

	producerMap = map[string]NewProducerFunc{
		"confluent-kafka": NewConfluentKafka,
		"plugin":          NewProducerPlugin,
	}
)

func init() {
	validateConfig()
}

func validateConfig() {
	var missingVariables string
	if BootstrapServers == "" {
		logger.Warn("Missing Env Config: 'PRODUCER_BOOTSTRAP_SERVERS', can't use Confluent Kafka Producer")
	}

	if ProducerType == "" {
		missingVariables = fmt.Sprintf("%s PRODUCER_TYPE", missingVariables)
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
