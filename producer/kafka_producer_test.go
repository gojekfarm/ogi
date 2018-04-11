package ogiproducer

import (
	"net/http"
	"testing"

	"github.com/abhishekkr/gol/golerror"
	"github.com/bouk/monkey"
	kafka "github.com/confluentinc/confluent-kafka-go/kafka"
	newrelic "github.com/newrelic/go-agent"
	"github.com/stretchr/testify/assert"

	"github.com/gojekfarm/ogi/instrumentation"
	"github.com/gojekfarm/ogi/logger"
)

func TestKafkaNewProducerSuccess(t *testing.T) {
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

func TestKafkaNewProducerFailure(t *testing.T) {
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

func TestKafkaNewProducerFailAtValidateConfig(t *testing.T) {
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

func TestKafkaClose(t *testing.T) {
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

func TestKafkaGetMetadataSuccessForEmptyMetadata(t *testing.T) {
	setTestConfig()
	k := &Kafka{}

	var nr, nrEnd, kguard *monkey.PatchGuard
	var kguardB bool
	nr = monkey.Patch(instrumentation.StartTransaction, func(string, http.ResponseWriter, *http.Request) newrelic.Transaction {
		nr.Unpatch()
		defer nr.Restore()
		return nil
	})
	nrEnd = monkey.Patch(instrumentation.EndTransaction, func(*newrelic.Transaction) {
		nrEnd.Unpatch()
		defer nrEnd.Restore()
		return
	})
	kguard = monkey.Patch((*kafka.Producer).GetMetadata, func(*kafka.Producer, *string, bool, int) (*kafka.Metadata, error) {
		kguard.Unpatch()
		defer kguard.Restore()
		kguardB = true
		return &kafka.Metadata{}, nil
	})

	k.GetMetadata()
	assert.True(t, kguardB)
	assert.Equal(t, map[string]int{}, k.PartitionCounts)
}

func TestKafkaGetMetadataSuccessForSomeMetadata(t *testing.T) {
	setTestConfig()
	k := &Kafka{}

	var nr, nrEnd, kguard *monkey.PatchGuard
	var kguardB bool
	nr = monkey.Patch(instrumentation.StartTransaction, func(string, http.ResponseWriter, *http.Request) newrelic.Transaction {
		nr.Unpatch()
		defer nr.Restore()
		return nil
	})
	nrEnd = monkey.Patch(instrumentation.EndTransaction, func(*newrelic.Transaction) {
		nrEnd.Unpatch()
		defer nrEnd.Restore()
		return
	})
	kguard = monkey.Patch((*kafka.Producer).GetMetadata, func(*kafka.Producer, *string, bool, int) (*kafka.Metadata, error) {
		kguard.Unpatch()
		defer kguard.Restore()
		kguardB = true
		return &kafka.Metadata{
			Topics: map[string]kafka.TopicMetadata{
				"some":  kafka.TopicMetadata{},
				"other": kafka.TopicMetadata{},
			},
		}, nil
	})

	k.GetMetadata()
	assert.True(t, kguardB)
	assert.Equal(t, 2, len(k.PartitionCounts))
}

func TestKafkaGetMetadataFailure(t *testing.T) {
	setTestConfig()
	k := &Kafka{}

	var nr, nrEnd, kguard *monkey.PatchGuard
	var kguardB bool
	nr = monkey.Patch(instrumentation.StartTransaction, func(string, http.ResponseWriter, *http.Request) newrelic.Transaction {
		nr.Unpatch()
		defer nr.Restore()
		return nil
	})
	nrEnd = monkey.Patch(instrumentation.EndTransaction, func(*newrelic.Transaction) {
		nrEnd.Unpatch()
		defer nrEnd.Restore()
		return
	})
	kguard = monkey.Patch((*kafka.Producer).GetMetadata, func(*kafka.Producer, *string, bool, int) (*kafka.Metadata, error) {
		kguard.Unpatch()
		defer kguard.Restore()
		kguardB = true
		return &kafka.Metadata{
			Topics: map[string]kafka.TopicMetadata{
				"some":  kafka.TopicMetadata{},
				"other": kafka.TopicMetadata{},
			},
		}, golerror.Error(123, "this is error")
	})

	k.GetMetadata()
	assert.True(t, kguardB)
	assert.Equal(t, 0, len(k.PartitionCounts))
}

func TestKafkaGetPartitionNumberForValidPartitionCount(t *testing.T) {
	setTestConfig()
	k := &Kafka{
		PartitionCounts: map[string]int{"some": 10, "other": 7},
	}

	var nr, nrEnd *monkey.PatchGuard
	nr = monkey.Patch(instrumentation.StartTransaction, func(string, http.ResponseWriter, *http.Request) newrelic.Transaction {
		nr.Unpatch()
		defer nr.Restore()
		return nil
	})
	nrEnd = monkey.Patch(instrumentation.EndTransaction, func(*newrelic.Transaction) {
		nrEnd.Unpatch()
		defer nrEnd.Restore()
		return
	})

	partitionNum := k.GetPartitionNumber("some", "what")
	assert.Equal(t, int(partitionNum), 4)
	partitionNum = k.GetPartitionNumber("other", "ever")
	assert.Equal(t, int(partitionNum), 6)
}

func TestKafkaGetPartitionNumberForMissingPartitionCount(t *testing.T) {
	setTestConfig()
	k := &Kafka{}

	var nr, nrEnd, kguard *monkey.PatchGuard
	var kguardB bool
	nr = monkey.Patch(instrumentation.StartTransaction, func(string, http.ResponseWriter, *http.Request) newrelic.Transaction {
		nr.Unpatch()
		defer nr.Restore()
		return nil
	})
	nrEnd = monkey.Patch(instrumentation.EndTransaction, func(*newrelic.Transaction) {
		nrEnd.Unpatch()
		defer nrEnd.Restore()
		return
	})
	kguard = monkey.Patch((*Kafka).GetMetadata, func(_k *Kafka) {
		kguard.Unpatch()
		defer kguard.Restore()
		kguardB = true
		_k.PartitionCounts = map[string]int{"some": 10, "other": 7}
	})

	partitionNum := k.GetPartitionNumber("some", "what")
	assert.Equal(t, int(partitionNum), 4)
	partitionNum = k.GetPartitionNumber("other", "ever")
	assert.Equal(t, int(partitionNum), 6)
	assert.True(t, kguardB)
}

func TestKafkaGetPartitionNumberForZeroPartitionCount(t *testing.T) {
	setTestConfig()
	k := &Kafka{
		PartitionCounts: map[string]int{"other": 0},
	}

	var nr, nrEnd, kguard *monkey.PatchGuard
	nr = monkey.Patch(instrumentation.StartTransaction, func(string, http.ResponseWriter, *http.Request) newrelic.Transaction {
		nr.Unpatch()
		defer nr.Restore()
		return nil
	})
	nrEnd = monkey.Patch(instrumentation.EndTransaction, func(*newrelic.Transaction) {
		nrEnd.Unpatch()
		defer nrEnd.Restore()
		return
	})
	kguard = monkey.Patch((*Kafka).GetMetadata, func(_k *Kafka) {
		kguard.Unpatch()
		defer kguard.Restore()
		_k.PartitionCounts = map[string]int{"other": 0}
	})

	partitionNum := k.GetPartitionNumber("other", "ever")
	assert.Equal(t, int(partitionNum), -1) //checks in kafka.PartitionAny is returned
}

func xTestKafkaProduceMessage(t *testing.T) {
	//TBD
	assert.Equal(t, "wip", "wip")
}
