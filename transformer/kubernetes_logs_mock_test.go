package ogitransformer

import "github.com/stretchr/testify/mock"

type MockKubernetesKafkaLog struct {
	mock.Mock
}

func (k *MockKubernetesKafkaLog) Transform(msg string) error {
	k.Mock.Called()
	return nil
}

func NewMockKafkaLog() Transformer {
	return &MockKubernetesKafkaLog{}
}
