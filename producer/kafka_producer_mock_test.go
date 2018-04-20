package ogiproducer

import (
	logger "github.com/gojekfarm/ogi/logger"
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

func (k *MockProducer) Produce(topic string, message []byte, messageKey string) {
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
