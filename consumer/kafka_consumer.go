package ogiconsumer

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	logger "github.com/gojekfarm/kafka-ogi/logger"
	ogiproducer "github.com/gojekfarm/kafka-ogi/producer"
	ogitransformer "github.com/gojekfarm/kafka-ogi/transformer"
)

type Kafka struct {
	ConfigMap kafka.ConfigMap
	Consumer  *kafka.Consumer
}

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

func (k *Kafka) SubscribeTopics(topics []string) {
	err := k.Consumer.SubscribeTopics(topics, nil)
	failIfError(err)
	logger.Debug(k.Consumer)
}

func (k *Kafka) EventHandler() {
	producer := ogiproducer.NewProducer()
	defer ogiproducer.CloseProducer(producer)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigchan:
			logger.Errorf("Caught signal %v: terminating\n", sig)
			break

		case ev := <-k.Consumer.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				logger.Infof("%% %v\n", e)
				k.Consumer.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				logger.Infof("%% %v\n", e)
				k.Consumer.Unassign()
			case *kafka.Message:
				////fmt.Printf("%% Message on %s:\n%s\n",
				////	e.TopicPartition, string(e.Value))
				ogitransformer.Transform(producer, string(e.Value))

			case kafka.PartitionEOF:
				logger.Infof("%% Reached %v\n", e)
			case kafka.Error:
				logger.Errorf("%% Error: %v\n", e)
				break
			}
		}
	}
}

func (k *Kafka) Close() {
	k.Consumer.Close()
}
