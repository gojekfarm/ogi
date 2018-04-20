package ogiconsumer

import (
	logger "github.com/gojekfarm/ogi/logger"
	"github.com/stretchr/testify/mock"
)

type MockConsumer struct {
	mock.Mock
}

func (k *MockConsumer) Consume() {
	k.Mock.Called()
	return
}

func NewMockConsumer() Consumer {
	var k MockConsumer
	return &k
}

func setTestConfig() {
	ConsumerType = "mock"
	consumerMap["mock"] = NewMockConsumer

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
