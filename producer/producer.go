package ogiproducer

import (
	"fmt"
	"hash/crc32"

	"github.com/abhishekkr/gol/golenv"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gojekfarm/kafka-ogi/instrumentation"
	"github.com/gojekfarm/kafka-ogi/logger"
)

var (
	BootstrapServers = golenv.OverrideIfEnv("PRODUCER_BOOTSTRAP_SERVERS", "")
	PartitionCounts  map[string]int
)

func validateConfig() {
	var missingVariables string
	if BootstrapServers == "" {
		missingVariables = fmt.Sprintf("%s CONSUMER_BOOTSTRAP_SERVERS", missingVariables)
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

func NewProducer() *kafka.Producer {
	validateConfig()
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": BootstrapServers,
	})

	failIfError(err)
	return producer
}

func CloseProducer(producer *kafka.Producer) {
	producer.Close()
}

func getMetadata(producer *kafka.Producer) {
	txn := instrumentation.StartTransaction("getMetdata_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	metadata, err := producer.GetMetadata(nil, true, 3000)
	if err != nil {
		return
	}
	PartitionCounts = make(map[string]int, len(metadata.Topics))
	for topic, _ := range metadata.Topics {
		PartitionCounts[topic] = len(metadata.Topics[topic].Partitions)
	}
}

func getPartitionNumber(producer *kafka.Producer, topic string, messageKey string) (partitionNumber int32) {
	txn := instrumentation.StartTransaction("getPartitionNumber_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	if len(PartitionCounts) == 0 || PartitionCounts[topic] == 0 {
		getMetadata(producer)
	}

	if PartitionCounts[topic] != 0 {
		partitionCount := PartitionCounts[topic]
		partitionChecksum := crc32.ChecksumIEEE([]byte(messageKey))
		partitionNumber = int32(partitionChecksum % uint32(partitionCount))
	} else {
		partitionNumber = kafka.PartitionAny
	}
	return
}

func produceMessage(producer *kafka.Producer, topic string, message []byte, partitionNumber int32) {
	txn := instrumentation.StartTransaction("produceMessage_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)
	doneChan := make(chan bool)

	go func() {
		defer close(doneChan)
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					logger.Errorf("Delivery failed for topic '%s[%d]' : %v\n",
						topic, partitionNumber, m.TopicPartition.Error)
				} else {
					logger.Infof("Delivered message to topic '%s' [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
				return

			default:
				logger.Infof("Ignored event: %s\n", ev)
			}
		}
	}()

	producer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &topic, Partition: partitionNumber},
		Value: message}

	// wait for delivery report goroutine to finish
	_ = <-doneChan
}

func Produce(producer *kafka.Producer, topic string, message []byte, message_key string) {
	txn := instrumentation.StartTransaction("produce_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)
	partitionNumber := getPartitionNumber(producer, topic, message_key)
	produceMessage(producer, topic, message, partitionNumber)
	logger.Infof("topic '%s' message-key '%s' is assigned to '%d' partition", topic, message_key, partitionNumber)
}
