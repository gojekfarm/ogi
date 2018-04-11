package ogitransformer

import (
	"net/http"
	"testing"

	"github.com/abhishekkr/gol/golerror"
	"github.com/bouk/monkey"
	newrelic "github.com/newrelic/go-agent"
	"github.com/stretchr/testify/assert"

	instrumentation "github.com/gojekfarm/ogi/instrumentation"
	logger "github.com/gojekfarm/ogi/logger"
	ogiproducer "github.com/gojekfarm/ogi/producer"
)

func setTestConfig() {
	KafkaTopicLabel = "app"
	TransformerType = "kubernetes-kafka-log"

	transformerMap = map[string]NewLogTransformer{
		"kubernetes-kafka-log": NewKubernetesKafkaLog,
		"mock":                 NewMockKafkaLog,
	}
}

func unsetTestConfig() {
	KafkaTopicLabel = ""
	TransformerType = ""

	transformerMap = map[string]NewLogTransformer{}
}

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

func TestTransformSuccess(t *testing.T) {
	setTestConfig()

	var nr, nrEnd, guard, logrguard *monkey.PatchGuard
	var nrB, nrEndB, guardB, logrguardB bool
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
	guard = monkey.Patch((*KubernetesKafkaLog).Transform, func(*KubernetesKafkaLog, string, ogiproducer.Producer) error {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		return nil
	})
	logrguard = monkey.Patch(logger.Warn, func(p ...interface{}) {
		logrguard.Unpatch()
		defer logrguard.Restore()
		logrguardB = true
		return
	})

	Transform(&ogiproducer.Kafka{}, "{}")
	assert.True(t, nrB)
	assert.True(t, nrEndB)
	assert.True(t, guardB)
	assert.False(t, logrguardB)
}

func TestTransformFailure(t *testing.T) {
	setTestConfig()

	var nr, nrEnd, guard, logrguard *monkey.PatchGuard
	var guardB, logrguardB bool
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
	guard = monkey.Patch((*KubernetesKafkaLog).Transform, func(*KubernetesKafkaLog, string, ogiproducer.Producer) error {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		return golerror.Error(123, "this is error")
	})
	logrguard = monkey.Patch(logger.Warn, func(p ...interface{}) {
		logrguard.Unpatch()
		defer logrguard.Restore()
		logrguardB = true
		return
	})

	Transform(&ogiproducer.Kafka{}, "{}")
	assert.True(t, guardB)
	assert.True(t, logrguardB)
}
