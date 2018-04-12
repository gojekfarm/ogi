package ogiproducer

import (
	"hash/crc32"

	kafka "github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/gojekfarm/ogi/instrumentation"
	logger "github.com/gojekfarm/ogi/logger"
)

type Kafka struct {
	ConfigMap       kafka.ConfigMap
	Producer        *kafka.Producer
	PartitionCounts map[string]int
}

func (k *Kafka) NewProducer() {
	var err error
	validateConfig()
	k.Producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": BootstrapServers,
	})

	if err != nil {
		logger.Fatal(err)
	}
	return
}

func (k *Kafka) Close() {
	k.Producer.Close()
}

func (k *Kafka) GetMetadata() {
	txn := instrumentation.StartTransaction("getMetdata_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	metadata, err := k.Producer.GetMetadata(nil, true, 3000)
	if err != nil {
		return
	}
	k.PartitionCounts = make(map[string]int, len(metadata.Topics))
	for topic, _ := range metadata.Topics {
		k.PartitionCounts[topic] = len(metadata.Topics[topic].Partitions)
	}
}

func (k *Kafka) GetPartitionNumber(topic string, messageKey string) (partitionNumber int32) {
	txn := instrumentation.StartTransaction("getPartitionNumber_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	if len(k.PartitionCounts) == 0 || k.PartitionCounts[topic] == 0 {
		k.GetMetadata()
	}

	if k.PartitionCounts[topic] != 0 {
		partitionCount := k.PartitionCounts[topic]
		partitionChecksum := crc32.ChecksumIEEE([]byte(messageKey))
		partitionNumber = int32(partitionChecksum % uint32(partitionCount))
	} else {
		partitionNumber = kafka.PartitionAny
	}
	return
}

func (k *Kafka) ProduceMessage(topic string, message []byte, partitionNumber int32) {
	txn := instrumentation.StartTransaction("produceMessage_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)
	doneChan := make(chan bool)

	go func() {
		defer close(doneChan)
		for e := range k.Producer.Events() {
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

	k.Producer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &topic, Partition: partitionNumber},
		Value: message}

	// wait for delivery report goroutine to finish
	_ = <-doneChan
}

func NewConfluentKafka() Producer {
	k := &Kafka{}
	k.NewProducer()
	return k
}
