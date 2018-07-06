package ogiproducer

import (
	"testing"

	"github.com/abhishekkr/gol/golenv"
	logger "github.com/gojektech/ogi/logger"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func TestNewProducerPluginOnSuccess(t *testing.T) {
	ProducerType = "plugin"
	ProducerPluginPath = golenv.OverrideIfEnv("PRODUCER_PLUGIN_PATH", "")
	var gLoggerFatalf *monkey.PatchGuard

	gLoggerFatalf = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		gLoggerFatalf.Unpatch()
		defer gLoggerFatalf.Restore()
		panic("mocked")
	})

	assert.NotPanics(t, func() { NewProducerPlugin() })
	ProducerType = golenv.OverrideIfEnv("PRODUCER_TYPE", "confluent-kafka")
}

func TestProducerPluginOnPluginFuncMissing(t *testing.T) {
	ProducerType = "plugin"
	ProducerPluginPath = golenv.OverrideIfEnv("PRODUCER_BAD_PLUGIN_PATH", "")
	var gLoggerFatalf *monkey.PatchGuard

	gLoggerFatalf = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		gLoggerFatalf.Unpatch()
		defer gLoggerFatalf.Restore()
		panic("mocked")
	})

	assert.Panics(t, func() { NewProducerPlugin() })
	ProducerPluginPath = golenv.OverrideIfEnv("PRODUCER_PLUGIN_PATH", "")
	ProducerType = golenv.OverrideIfEnv("PRODUCER_TYPE", "confluent-kafka")
}

func TestProducerPluginOnPluginFileMissing(t *testing.T) {
	ProducerType = "plugin"
	ProducerPluginPath = ""
	var gLoggerFatalf *monkey.PatchGuard

	gLoggerFatalf = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		gLoggerFatalf.Unpatch()
		defer gLoggerFatalf.Restore()
		panic("mocked")
	})

	assert.Panics(t, func() { NewProducerPlugin() })
	ProducerPluginPath = golenv.OverrideIfEnv("PRODUCER_PLUGIN_PATH", "")
	ProducerType = golenv.OverrideIfEnv("PRODUCER_TYPE", "confluent-kafka")
}
