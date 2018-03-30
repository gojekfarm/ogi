package ogitransformer

import (
	"encoding/json"
	"fmt"

	"github.com/abhishekkr/gol/golenv"

	"github.com/gojekfarm/kafka-ogi/instrumentation"

	"github.com/gojekfarm/kafka-ogi/logger"
	ogiproducer "github.com/gojekfarm/kafka-ogi/producer"
)

type KafkaLog struct {
	Message    string     `json:"message"`
	Stream     string     `json:"stream"`
	LogLine    string     `json:"log"`
	Docker     Docker     `json:"docker"`
	Kubernetes Kubernetes `json:"kubernetes"`
	MessageKey string     `json:"message_key"`
}

type Docker struct {
	ContainerId string `json:"container_id"`
}

type Kubernetes struct {
	ContainerName string            `json:"container_name"`
	NamespaceName string            `json:"namespace_name"`
	PodName       string            `json:"pod_name"`
	PodId         string            `json:"pod_id"`
	Labels        map[string]string `json:"labels"`
	Host          string            `json:"host"`
	MasterUrl     string            `json:"master_url"`
}

var (
	KafkaTopicLabel = golenv.OverrideIfEnv("PRODUCER_KAFKA_TOPIC_LABEL", "app")
)

func validateConfig() {
	var missingVariables string
	if KafkaTopicLabel == "" {
		missingVariables = fmt.Sprintf("%s PRODUCER_KAFKA_TOPIC_LABEL", missingVariables)
	}

	if missingVariables != "" {
		logger.Fatalf("Missing Env Config:%s", missingVariables)
	}
}

func Transform(producer ogiproducer.Producer, msg string) {
	txn := instrumentation.StartTransaction("transform_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)
	msgBytes := []byte(msg)

	var kafkaLog KafkaLog
	err := json.Unmarshal(msgBytes, &kafkaLog)

	if err != nil {
		logger.Warnf("failed to parse", msg)
	}

	topic := kafkaLog.Kubernetes.Labels[KafkaTopicLabel]
	kafkaLog.MessageKey = kafkaLog.Kubernetes.PodName

	msgWithKey, err := json.Marshal(kafkaLog)

	if topic != "" {
		ogiproducer.Produce(producer, topic, msgWithKey, kafkaLog.MessageKey)
	} else {
		logger.Warnf("correct target topic id is missing for",
			kafkaLog)
	}
}
