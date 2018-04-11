package ogiconsumer

import (
	"net/http"
	"testing"

	"github.com/abhishekkr/gol/golerror"
	"github.com/bouk/monkey"
	newrelic "github.com/newrelic/go-agent"
	"github.com/stretchr/testify/assert"

	instrumentation "github.com/gojekfarm/ogi/instrumentation"
	logger "github.com/gojekfarm/ogi/logger"
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

func TestFailIfError(t *testing.T) {
	var guard *monkey.PatchGuard
	guard = monkey.Patch(logger.Fatal, func(p ...interface{}) {
		guard.Unpatch()
		defer guard.Restore()

		panic("mocked")
		return
	})

	var thisErr error
	assert.NotPanics(t, func() { failIfError(thisErr) })
	thisErr = golerror.Error(123, "this is an error")
	assert.Panicsf(t, func() { failIfError(thisErr) }, "mocked")
}

func TestSubscribeForValidTopic(t *testing.T) {
	setTestConfig()
	mc := new(MockConsumer)
	mc.On("SubscribeTopics").Return()
	mc.On("EventHandler").Return()
	subscribe(mc)
}

func TestSubscribeForNoTopic(t *testing.T) {
	var guard *monkey.PatchGuard
	guard = monkey.Patch(logger.Fatal, func(p ...interface{}) {
		guard.Unpatch()
		defer guard.Restore()

		panic("mocked")
		return
	})
	unsetTestConfig()
	mc := new(MockConsumer)

	assert.Panicsf(t, func() { subscribe(mc) }, "mocked")
}

func TestConsume(t *testing.T) {
	var nr, nrEnd, vc, s *monkey.PatchGuard
	var nrB, nrEndB, vcB, sB bool
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
	vc = monkey.Patch(validateConfig, func() {
		vc.Unpatch()
		defer vc.Restore()
		vcB = true
		return
	})
	s = monkey.Patch(subscribe, func(Consumer) {
		s.Unpatch()
		defer s.Restore()
		sB = true
		return
	})
	setTestConfig()
	Consume()
	assert.Equal(t, nrB, true)
	assert.Equal(t, nrEndB, true)
	assert.Equal(t, vcB, true)
	assert.Equal(t, sB, true)
}
