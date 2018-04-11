package ogiconsumer

import (
	logger "github.com/gojekfarm/ogi/logger"
	"github.com/stretchr/testify/mock"
)

type MockConsumer struct {
	mock.Mock
}

func (k *MockConsumer) Configure() {
	k.Mock.Called()
	return
}

func (k *MockConsumer) NewConsumer() {
	k.Mock.Called()
	return
}

func (k *MockConsumer) SubscribeTopics() {
	k.Mock.Called()
	return
}

func (k *MockConsumer) EventHandler() {
	k.Mock.Called()
	return
}

func (k *MockConsumer) Close() {
	k.Mock.Called()
	return
}

func setTestConfig() {
	KafkaTopics = "bulk-topic"
	BootstrapServers = "my-kafaka.server"
	GroupId = "ogi-group"
	SessionTimeoutMs = "6000"
	GoEventsChannelEnable = "true"
	GoEventsChannelSize = "1000"
	GoApplicationRebalanceEnable = "true"
	logger.SetupLogger()
}

func unsetTestConfig() {
	KafkaTopics = ""
	BootstrapServers = ""
	GroupId = ""
	SessionTimeoutMs = ""
	GoEventsChannelEnable = ""
	GoEventsChannelSize = ""
	GoApplicationRebalanceEnable = ""
}
