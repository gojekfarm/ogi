package ogiconsumer

import (
	"fmt"

	"github.com/abhishekkr/gol/golenv"

	instrumentation "github.com/gojekfarm/ogi/instrumentation"
	logger "github.com/gojekfarm/ogi/logger"
)

type Consumer interface {
	Configure()
	NewConsumer()
	SubscribeTopics()
	EventHandler()
	Close()
}

type NewConsumerFunc func() Consumer

var (
	BootstrapServers             = golenv.OverrideIfEnv("CONSUMER_BOOTSTRAP_SERVERS", "")
	GroupId                      = golenv.OverrideIfEnv("CONSUMER_GROUP_ID", "")
	SessionTimeoutMs             = golenv.OverrideIfEnv("CONSUMER_SESSION_TIMEOUT_MS", "6000")
	GoEventsChannelEnable        = golenv.OverrideIfEnv("CONSUMER_GOEVENTS_CHANNEL_ENABLE", "true")
	GoEventsChannelSize          = golenv.OverrideIfEnv("CONSUMER_GOEVENTS_CHANNEL_SIZE", "1000")
	GoApplicationRebalanceEnable = golenv.OverrideIfEnv("CONSUMER_GO_APPLICATION_REBALANCE_ENABLE", "true")
	ConsumerType                 = golenv.OverrideIfEnv("CONSUMER_TYPE", "confluent-kafka")

	consumerMap = map[string]NewConsumerFunc{
		"confluent-kafka": NewConfluentKafka,
		"plugin":          NewConsumerPlugin,
	}
)

func validateConfig() {
	var missingVariables string
	if KafkaTopics == "" {
		logger.Warn("Missing Env Config: 'CONSUMER_KAFKA_TOPICS, can't use Kafka based transformers")
	}
	if BootstrapServers == "" {
		missingVariables = fmt.Sprintf("%s CONSUMER_BOOTSTRAP_SERVERS", missingVariables)
	}
	if GroupId == "" {
		missingVariables = fmt.Sprintf("%s CONSUMER_GROUP_ID", missingVariables)
	}

	if missingVariables != "" {
		logger.Fatalf("Missing Env Config:%s", missingVariables)
	}
}

func failIfError(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func subscribe(consumer Consumer) {
	txn := instrumentation.StartTransaction("subscribe_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)
	consumer.SubscribeTopics()

	consumer.EventHandler()
}

func Consume() {
	txn := instrumentation.StartTransaction("consume_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)
	validateConfig()

	consumer := consumerMap[ConsumerType]()
	consumer.Configure()
	consumer.NewConsumer()

	subscribe(consumer)

	logger.Infof("Closing consumer\n")
	consumer.Close()
}
