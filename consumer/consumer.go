package ogiconsumer

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/abhishekkr/gol/golenv"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gojekfarm/kafka-ogi/logger"

	instrumentation "github.com/gojekfarm/kafka-ogi/instrumentation"
	ogiproducer "github.com/gojekfarm/kafka-ogi/producer"
	ogitransformer "github.com/gojekfarm/kafka-ogi/transformer"
)

var (
	KafkaTopics                  = golenv.OverrideIfEnv("CONSUMER_KAFKA_TOPICS", "")
	BootstrapServers             = golenv.OverrideIfEnv("CONSUMER_BOOTSTRAP_SERVERS", "")
	GroupId                      = golenv.OverrideIfEnv("CONSUMER_GROUP_ID", "")
	SessionTimeoutMs             = golenv.OverrideIfEnv("CONSUMER_SESSION_TIMEOUT_MS", "6000")
	GoEventsChannelEnable        = golenv.OverrideIfEnv("CONSUMER_GOEVENTS_CHANNEL_ENABLE", "true")
	GoEventsChannelSize          = golenv.OverrideIfEnv("CONSUMER_GOEVENTS_CHANNEL_SIZE", "1000")
	GoApplicationRebalanceEnable = golenv.OverrideIfEnv("CONSUMER_GO_APPLICATION_REBALANCE_ENABLE", "true")
)

func validateConfig() {
	var missingVariables string
	if KafkaTopics == "" {
		missingVariables = fmt.Sprintf("%s CONSUMER_KAFKA_TOPICS", missingVariables)
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

func newConsumer() *kafka.Consumer {
	var err error
	sessionTimeoutMs, err := strconv.Atoi(SessionTimeoutMs)
	failIfError(err)
	goEventsChannelSize, err := strconv.Atoi(GoEventsChannelSize)
	failIfError(err)
	goEventsChannelEnable, err := strconv.ParseBool(GoEventsChannelEnable)
	failIfError(err)
	goApplicationRebalanceEnable, err := strconv.ParseBool(GoApplicationRebalanceEnable)
	failIfError(err)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               BootstrapServers,
		"group.id":                        GroupId,
		"session.timeout.ms":              sessionTimeoutMs,
		"go.events.channel.size":          goEventsChannelSize,
		"go.events.channel.enable":        goEventsChannelEnable,
		"go.application.rebalance.enable": goApplicationRebalanceEnable,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"}})

	if err != nil {
		logger.Errorf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	return consumer
}

func subscribe(consumer *kafka.Consumer) {
	txn := instrumentation.StartTransaction("subscribe_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	topics := strings.Split(KafkaTopics, ",")
	logger.Debug(topics)
	err := consumer.SubscribeTopics(topics, nil)
	failIfError(err)
	logger.Debug(consumer)

	producer := ogiproducer.NewProducer()
	defer ogiproducer.CloseProducer(producer)

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			logger.Errorf("Caught signal %v: terminating\n", sig)
			run = false

		case ev := <-consumer.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				logger.Infof("%% %v\n", e)
				consumer.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				logger.Infof("%% %v\n", e)
				consumer.Unassign()
			case *kafka.Message:
				////fmt.Printf("%% Message on %s:\n%s\n",
				////	e.TopicPartition, string(e.Value))
				ogitransformer.Transform(producer, string(e.Value))

			case kafka.PartitionEOF:
				logger.Infof("%% Reached %v\n", e)
			case kafka.Error:
				logger.Errorf("%% Error: %v\n", e)
				run = false
			}
		}
	}
}

func Consume() {
	txn := instrumentation.StartTransaction("consume_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)
	validateConfig()
	consumer := newConsumer()

	subscribe(consumer)

	logger.Infof("Closing consumer\n")
	consumer.Close()
}
