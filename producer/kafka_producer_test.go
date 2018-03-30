package ogiproducer

import (
	"testing"

	"github.com/abhishekkr/gol/golerror"
	"github.com/bouk/monkey"
	kafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"

	logger "github.com/gojekfarm/kafka-ogi/logger"
)

func TestNewProducerSuccess(t *testing.T) {
	setTestConfig()
	k := &Kafka{}

	var guard, kafkaGuard *monkey.PatchGuard
	var guardB, kafkaGuardB bool
	guard = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		panic("mocked")
	})
	kafkaGuard = monkey.Patch(kafka.NewProducer, func(*kafka.ConfigMap) (*kafka.Producer, error) {
		kafkaGuard.Unpatch()
		defer kafkaGuard.Restore()
		kafkaGuardB = true
		return &kafka.Producer{}, nil
	})

	assert.NotPanics(t, func() { k.NewProducer() }, "mocked")
	assert.False(t, guardB)
	assert.True(t, kafkaGuardB)
}

func TestNewProducerFailure(t *testing.T) {
	setTestConfig()
	k := &Kafka{}

	var guard, kafkaGuard *monkey.PatchGuard
	var guardB, kafkaGuardB bool
	guard = monkey.Patch(logger.Fatal, func(p ...interface{}) {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		panic("mocked")
	})
	kafkaGuard = monkey.Patch(kafka.NewProducer, func(*kafka.ConfigMap) (*kafka.Producer, error) {
		kafkaGuard.Unpatch()
		defer kafkaGuard.Restore()
		kafkaGuardB = true
		return &kafka.Producer{}, golerror.Error(123, "this is an error")
	})

	assert.Panicsf(t, func() { k.NewProducer() }, "mocked")
	assert.True(t, guardB)
	assert.True(t, kafkaGuardB)
}

func TestNewProducerFailAtValidateConfig(t *testing.T) {
	unsetTestConfig()
	k := &Kafka{}

	var guard *monkey.PatchGuard
	var guardB bool
	guard = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		panic("mocked")
	})

	assert.Panicsf(t, func() { k.NewProducer() }, "mocked")
	assert.True(t, guardB)
}

func TestClose(t *testing.T) {
	setTestConfig()
	k := &Kafka{}

	var guard *monkey.PatchGuard
	var guardB bool
	guard = monkey.Patch((*kafka.Producer).Close, func(*kafka.Producer) {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		return
	})

	k.Close()
	assert.True(t, guardB)
}

func TestGetMetadataSuccess(t *testing.T) {
	//TBD
	assert.Equal(t, "wip", "tbd")
}

func TestGetMetadataFailure(t *testing.T) {
	//TBD
	assert.Equal(t, "wip", "tbd")
}

func TestGetPartitionNumberForValidPartitionCount(t *testing.T) {
	//TBD
	assert.Equal(t, "wip", "tbd")
}

func TestGetPartitionNumberForMissingPartitionCount(t *testing.T) {
	//TBD
	assert.Equal(t, "wip", "tbd")
}

func TestProduceMessage(t *testing.T) {
	//TBD
	assert.Equal(t, "wip", "tbd")
}
