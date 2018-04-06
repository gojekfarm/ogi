package ogitransformer

import (
	"github.com/stretchr/testify/mock"

	ogiproducer "github.com/gojekfarm/kafka-ogi/producer"
)

type MockKubernetesKafkaLog struct {
	mock.Mock
}

func (k *MockKubernetesKafkaLog) Transform(msg string, producer ogiproducer.Producer) error {
	k.Mock.Called()
	return nil
}

func NewMockKafkaLog() LogTransformer {
	return &MockKubernetesKafkaLog{}
}
