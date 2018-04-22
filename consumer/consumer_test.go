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

func TestConsume(t *testing.T) {
	var nr, nrEnd, mockGuard *monkey.PatchGuard
	var nrB, nrEndB bool
	mc := MockConsumer{}
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
	mockGuard = monkey.Patch(NewMockConsumer, func() Consumer {
		mockGuard.Unpatch()
		defer mockGuard.Restore()
		return &mc
	})

	mc.On("Consume").Return()
	setTestConfig()
	Consume()
	assert.Equal(t, nrB, true)
	assert.Equal(t, nrEndB, true)
	mc.Mock.AssertExpectations(t)
}
