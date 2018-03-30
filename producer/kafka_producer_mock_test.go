package ogiproducer

import (
	logger "github.com/gojekfarm/kafka-ogi/logger"
	"github.com/stretchr/testify/mock"
)

type MockProducer struct {
	mock.Mock
}

func (k *MockProducer) NewProducer() {
	k.Mock.Called()
	return
}

func (k *MockProducer) Close() {
	k.Mock.Called()
	return
}

func (k *MockProducer) GetMetadata() {
	k.Mock.Called()
	return
}

func (k *MockProducer) GetPartitionNumber(topic string, messageKey string) (partitionNumber int32) {
	k.Mock.Called()
	return
}

func (k *MockProducer) ProduceMessage(topic string, message []byte, partitionNumber int32) {
	k.Mock.Called()
	return
}

func setTestConfig() {
	BootstrapServers = "someserver"
	logger.SetupLogger()
}

func unsetTestConfig() {
	BootstrapServers = ""
}
