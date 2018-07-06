package ogiconsumer

import (
	"fmt"
	"testing"

	"github.com/abhishekkr/gol/golenv"
	logger "github.com/gojektech/ogi/logger"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func TestNewConsumerPluginOnSuccess(t *testing.T) {
	ConsumerType = "plugin"
	ConsumerPluginPath = golenv.OverrideIfEnv("CONSUMER_PLUGIN_PATH", "")
	var gLoggerFatalf *monkey.PatchGuard

	gLoggerFatalf = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		gLoggerFatalf.Unpatch()
		defer gLoggerFatalf.Restore()
		fmt.Printf(f, p...)
		panic("mocked")
	})

	assert.NotPanics(t, func() { NewConsumerPlugin() })
	ConsumerType = golenv.OverrideIfEnv("CONSUMER_TYPE", "confluent-kafka")
}

func TestConsumerPluginOnPluginFuncMissing(t *testing.T) {
	ConsumerType = "plugin"
	ConsumerPluginPath = golenv.OverrideIfEnv("CONSUMER_BAD_PLUGIN_PATH", "")
	var gLoggerFatalf *monkey.PatchGuard

	gLoggerFatalf = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		gLoggerFatalf.Unpatch()
		defer gLoggerFatalf.Restore()
		fmt.Printf(f, p...)
		panic("mocked")
	})

	assert.Panics(t, func() { NewConsumerPlugin() })
	ConsumerPluginPath = golenv.OverrideIfEnv("CONSUMER_PLUGIN_PATH", "")
	ConsumerType = golenv.OverrideIfEnv("CONSUMER_TYPE", "confluent-kafka")
}

func TestConsumerPluginOnPluginFileMissing(t *testing.T) {
	ConsumerType = "plugin"
	ConsumerPluginPath = ""
	var gLoggerFatalf *monkey.PatchGuard

	gLoggerFatalf = monkey.Patch(logger.Fatalf, func(f string, p ...interface{}) {
		gLoggerFatalf.Unpatch()
		defer gLoggerFatalf.Restore()
		fmt.Printf(f, p...)
		panic("mocked")
	})

	assert.Panics(t, func() { NewConsumerPlugin() })
	ConsumerPluginPath = golenv.OverrideIfEnv("CONSUMER_PLUGIN_PATH", "")
	ConsumerType = golenv.OverrideIfEnv("CONSUMER_TYPE", "confluent-kafka")
}
