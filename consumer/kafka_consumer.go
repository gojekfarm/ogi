package ogiconsumer

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/abhishekkr/gol/golenv"
	kafka "github.com/confluentinc/confluent-kafka-go/kafka"

	logger "github.com/gojektech/ogi/logger"
	ogitransformer "github.com/gojektech/ogi/transformer"
)

type Kafka struct {
	ConfigMap kafka.ConfigMap
	Consumer  *kafka.Consumer
}

var (
	KafkaTopics = golenv.OverrideIfEnv("CONSUMER_KAFKA_TOPICS", "")
)

func (k *Kafka) Configure() {
	sessionTimeoutMs, err := strconv.Atoi(SessionTimeoutMs)
	failIfError(err)
	goEventsChannelSize, err := strconv.Atoi(GoEventsChannelSize)
	failIfError(err)
	goEventsChannelEnable, err := strconv.ParseBool(GoEventsChannelEnable)
	failIfError(err)
	goApplicationRebalanceEnable, err := strconv.ParseBool(GoApplicationRebalanceEnable)
	failIfError(err)

	k.ConfigMap = kafka.ConfigMap{
		"bootstrap.servers":               BootstrapServers,
		"group.id":                        GroupId,
		"session.timeout.ms":              sessionTimeoutMs,
		"go.events.channel.size":          goEventsChannelSize,
		"go.events.channel.enable":        goEventsChannelEnable,
		"go.application.rebalance.enable": goApplicationRebalanceEnable,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}
}

func (k *Kafka) NewConsumer() {
	var err error
	k.Consumer, err = kafka.NewConsumer(&k.ConfigMap)

	if err != nil {
		logger.Fatalf("Failed to create consumer: %s\n", err)
	}
}

func (k *Kafka) SubscribeTopics() {
	logger.Debug(KafkaTopics)
	topics := strings.Split(KafkaTopics, ",")
	if len(topics) == 1 && topics[0] == "" {
		logger.Fatal("no topic provided to consume")
	}

	err := k.Consumer.SubscribeTopics(topics, nil)
	failIfError(err)
	logger.Debug(k.Consumer)
}

func (k *Kafka) EventHandler() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			logger.Errorf("Caught signal %v: terminating\n", sig)
			run = false

		case ev := <-k.Consumer.Events():
			switch e := ev.(type) {

			case kafka.AssignedPartitions:
				logger.Infof("%% %v\n", e)
				k.Consumer.Assign(e.Partitions)

			case kafka.RevokedPartitions:
				logger.Infof("%% %v\n", e)
				k.Consumer.Unassign()

			case *kafka.Message:
				go ogitransformer.Transform(e.Value)

			case kafka.PartitionEOF:
				logger.Infof("%% Reached %v\n", e)

			case kafka.Error:
				logger.Errorf("%% Error: %v\n", e)
				run = false
			}
		}
	}
}

func (k *Kafka) Close() {
	k.Consumer.Close()
}

func (k *Kafka) Consume() {
	k.Configure()
	k.NewConsumer()

	k.SubscribeTopics()

	k.EventHandler()

	logger.Infof("Closing consumer\n")
	k.Close()
}

func NewConfluentKafka() Consumer {
	if KafkaTopics == "" {
		logger.Fatalln("Missing Env Config: 'CONSUMER_KAFKA_TOPICS, can't use Kafka based transformers")
	}

	var k Kafka
	return &k
}
