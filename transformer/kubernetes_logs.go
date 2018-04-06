package ogitransformer

import (
	"encoding/json"
	"fmt"

	"github.com/abhishekkr/gol/golerror"

	ogiproducer "github.com/gojekfarm/kafka-ogi/producer"
)

type KubernetesKafkaLog struct {
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

func (kafkaLog *KubernetesKafkaLog) Transform(msg string, producer ogiproducer.Producer) (err error) {
	msgBytes := []byte(msg)

	if err = json.Unmarshal(msgBytes, &kafkaLog); err != nil {
		err = golerror.Error(123, "failed to parse")
		return
	}

	if kafkaLog.Kubernetes.Labels[KafkaTopicLabel] == "" {
		err = golerror.Error(123, fmt.Sprintf("correct target topic id '%s' is missing", KafkaTopicLabel))
		return
	}

	kafkaLog.MessageKey = kafkaLog.Kubernetes.PodName

	msgWithKey, err := json.Marshal(kafkaLog)

	ogiproducer.Produce(producer,
		kafkaLog.Kubernetes.Labels[KafkaTopicLabel],
		msgWithKey,
		kafkaLog.MessageKey)
	return
}

func NewKubernetesKafkaLog() LogTransformer {
	return &KubernetesKafkaLog{}
}
