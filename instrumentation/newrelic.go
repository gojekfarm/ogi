package instrumentation

import (
	"net/http"

	"github.com/abhishekkr/gol/golenv"
	newrelic "github.com/newrelic/go-agent"

	"github.com/gojekfarm/ogi/logger"
)

type NewrelicCtx struct {
	NewrelicApp newrelic.Application
}

var (
	instrumentationCtx *NewrelicCtx
)

func init() {
	config := newrelic.NewConfig(
		golenv.OverrideIfEnv("NEWRELIC_APP_NAME", "kafka_ogi"),
		golenv.OverrideIfEnv("NEWRELIC_LICENSE_KEY", "x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x"),
	)
	newrelicApp, err := newrelic.NewApplication(config)
	if err != nil {
		logger.Errorf("newrelic init failed: %s", err)
		return
	}
	instrumentationCtx = &NewrelicCtx{
		NewrelicApp: newrelicApp,
	}
	return
}

func StartTransaction(txn string, w http.ResponseWriter, r *http.Request) newrelic.Transaction {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			logger.Errorf("start transaction failed: %s", err)
		}
	}()

	if instrumentationCtx.NewrelicApp != nil {
		return instrumentationCtx.NewrelicApp.StartTransaction(txn, w, r)
	}
	logger.Errorln("newrelic is not initialized")
	return nil
}

func EndTransaction(transaction *newrelic.Transaction) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			logger.Errorf("end transaction failed: %s", err)
		}
	}()

	(*transaction).End()
}
