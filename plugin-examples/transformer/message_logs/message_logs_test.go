package main

import (
	"testing"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/mock"

	logger "github.com/gojekfarm/ogi/logger"
	ogiproducer "github.com/gojekfarm/ogi/producer"
)

type MockMessageLog struct {
	mock.Mock
}

func (msgLog *MockMessageLog) Transform(msg string, producer ogiproducer.Producer) (err error) {
	msgLog.Mock.Called()
	return
}

func TestMessageTransform(t *testing.T) {
	logger.SetupLogger()
	msgLog := MessageLog{}

	var guard *monkey.PatchGuard
	var guardB bool
	var guardMessage, guardKey, guardTopic string
	guard = monkey.Patch(ogiproducer.Produce, func(topic string, message []byte, message_key string) {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		guardMessage, guardKey, guardTopic = string(message), message_key, topic
		return
	})
	msgLog.Transform("whatever")
}
