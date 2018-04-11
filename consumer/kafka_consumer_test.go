package ogiconsumer

import (
	"testing"

	kafka "github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/abhishekkr/gol/golerror"
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"

	logger "github.com/gojekfarm/ogi/logger"
)

func TestKafkaConfigureWithValidConfig(t *testing.T) {
	setTestConfig()
	k := &Kafka{}

	var guard *monkey.PatchGuard
	var guardB bool
	guard = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		panic("mocked")
		return
	})

	assert.NotPanics(t, func() { k.Configure() }, "mocked")
	assert.False(t, guardB)
}

func TestKafkaConfigureTypeConvertsConfig(t *testing.T) {
	setTestConfig()
	k := &Kafka{}

	k.Configure()
	assert.Equal(t, "my-kafaka.server", k.ConfigMap["bootstrap.servers"])
	assert.Equal(t, "ogi-group", k.ConfigMap["group.id"])
	assert.Equal(t, 6000, k.ConfigMap["session.timeout.ms"])
	assert.Equal(t, 1000, k.ConfigMap["go.events.channel.size"])
	assert.Equal(t, true, k.ConfigMap["go.events.channel.enable"])
	assert.Equal(t, true, k.ConfigMap["go.application.rebalance.enable"])
}

func TestKafkaConfigureWithInvalidConfig(t *testing.T) {
	unsetTestConfig()
	k := &Kafka{}

	var guard *monkey.PatchGuard
	guard = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		guard.Unpatch()
		defer guard.Restore()
		panic("mocked")
	})

	assert.Panicsf(t, func() { k.Configure() }, "mocked")
}

func TestKafkaNewConsumerCallConsumerSuccess(t *testing.T) {
	setTestConfig()
	k := &Kafka{}

	var guard *monkey.PatchGuard
	var guardB bool
	guard = monkey.Patch(kafka.NewConsumer, func(*kafka.ConfigMap) (*kafka.Consumer, error) {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		return &kafka.Consumer{}, nil
	})

	k.NewConsumer()
	assert.True(t, guardB)
}

func TestKafkaNewConsumerCallConsumerFailure(t *testing.T) {
	unsetTestConfig()
	k := &Kafka{}

	var guard, failguard *monkey.PatchGuard
	var guardB bool
	guard = monkey.Patch(kafka.NewConsumer, func(*kafka.ConfigMap) (*kafka.Consumer, error) {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		return &kafka.Consumer{}, golerror.Error(123, "this is an error")
	})
	failguard = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		failguard.Unpatch()
		defer failguard.Restore()
		panic("mocked")
	})

	assert.Panicsf(t, func() { k.NewConsumer() }, "mocked")
	assert.True(t, guardB)
}

func TestKafkaSubscribeTopicsSuccess(t *testing.T) {
	setTestConfig()
	k := &Kafka{}

	var guard *monkey.PatchGuard
	var guardB bool
	guard = monkey.Patch((*kafka.Consumer).SubscribeTopics, func(*kafka.Consumer, []string, kafka.RebalanceCb) error {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		return nil
	})

	k.SubscribeTopics([]string{"bulk-topic"})
	assert.True(t, guardB)
}

func TestKafkaSubscribeTopicsFailure(t *testing.T) {
	unsetTestConfig()
	k := &Kafka{}

	var guard, failguard *monkey.PatchGuard
	var guardB bool
	guard = monkey.Patch((*kafka.Consumer).SubscribeTopics, func(*kafka.Consumer, []string, kafka.RebalanceCb) error {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		return golerror.Error(123, "this is an error")
	})
	failguard = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		failguard.Unpatch()
		defer failguard.Restore()
		panic("mocked")
	})

	assert.Panicsf(t, func() { k.SubscribeTopics([]string{"bulk-topic"}) }, "mocked")
	assert.True(t, guardB)
}

func TestKafkaEventHandler(t *testing.T) {
	//TBD
}

func TestKafkaClose(t *testing.T) {
	setTestConfig()
	k := &Kafka{}

	var guard *monkey.PatchGuard
	var guardB bool
	guard = monkey.Patch((*kafka.Consumer).Close, func(*kafka.Consumer) error {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		return nil
	})

	k.Close()
	assert.True(t, guardB)
}
