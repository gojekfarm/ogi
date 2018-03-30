package ogiproducer

import (
	"net/http"
	"testing"

	"github.com/bouk/monkey"
	newrelic "github.com/newrelic/go-agent"
	"github.com/stretchr/testify/assert"

	instrumentation "github.com/gojekfarm/kafka-ogi/instrumentation"
	logger "github.com/gojekfarm/kafka-ogi/logger"
)

func TestValdiateConfig(t *testing.T) {
	var guard *monkey.PatchGuard
	guard = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		guard.Unpatch()
		defer guard.Restore()

		panic("mocked")
		return
	})

	setTestConfig()
	assert.NotPanics(t, func() { validateConfig() })
	unsetTestConfig()
	assert.Panicsf(t, func() { validateConfig() }, "mocked")
}

func TestNewProducer(t *testing.T) {
	setTestConfig()

	var guard *monkey.PatchGuard
	var guardB bool
	guard = monkey.Patch((*Kafka).NewProducer, func(*Kafka) {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		return
	})

	NewProducer()
	assert.True(t, guardB)
}

func TestProduce(t *testing.T) {
	setTestConfig()

	var nr, nrEnd *monkey.PatchGuard
	var nrB, nrEndB bool
	nr = monkey.Patch(instrumentation.StartTransaction, func(string, http.ResponseWriter, *http.Request) newrelic.Transaction {
		nr.Unpatch()
		defer nr.Restore()
		nrB = true
		return nil
	})
	nrEnd = monkey.Patch(instrumentation.EndTransaction, func(*newrelic.Transaction) {
		nrEnd.Unpatch()
		defer nrEnd.Restore()
		nrEndB = true
		return
	})

	mp := &MockProducer{}
	mp.On("GetPartitionNumber").Return(1)
	mp.On("ProduceMessage").Return()
	Produce(mp, "topik", []byte{}, "key")
	assert.True(t, nrB)
	assert.True(t, nrEndB)
	mp.Mock.AssertExpectations(t)
}
