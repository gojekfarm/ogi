package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"github.com/abhishekkr/gol/golenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	instrumentation "github.com/gojekfarm/ogi/instrumentation"
	logger "github.com/gojekfarm/ogi/logger"
	ogiproducer "github.com/gojekfarm/ogi/producer"
)

type GCPStackdriverLogsConsumer struct {
}

type GCPLogEntry struct {
	InsertID    string               `json:"insert-id,omitempty"`
	Timestamp   time.Time            `json:"timestamp,omitempty"`
	Severity    string               `json:"severity,omitempty"`
	Labels      map[string]string    `json:"labels,omitempty"`
	HTTPRequest *logging.HTTPRequest `json:"http-request,omitempty"`
	Payload     interface{}          `json:"payload,omitempty"`
}

var (
	GoogleProjectId           = golenv.OverrideIfEnv("GOOGLE_PROJECT_ID", "my-project-100")
	GoogleLogTopic            = golenv.OverrideIfEnv("GOOGLE_LOG_TOPIC", "my-app")
	GoogleCloudCredentialFile = golenv.OverrideIfEnv("GOOGLE_CLOUD_CREDENTIAL_FILE", "/tmp/my-project-100-sa01.json")

	gcp *GCPStackdriverLogsConsumer

	severityLevels = map[int]string{
		0:   "default",
		100: "debug",
		200: "info",
		300: "notice",
		400: "warning",
		500: "error",
		600: "critical",
		700: "alert",
		800: "emergency",
	}
)

func init() {
	gcp = new(GCPStackdriverLogsConsumer)
}

func (gcp *GCPStackdriverLogsConsumer) produce(entry *logging.Entry) {
	severity := severityLevels[int(entry.Severity)]
	gcpLog := GCPLogEntry{
		InsertID:    entry.InsertID,
		Timestamp:   entry.Timestamp,
		Severity:    severity,
		Labels:      entry.Labels,
		HTTPRequest: entry.HTTPRequest,
		Payload:     entry.Payload,
	}
	logbytes, err := json.Marshal(gcpLog)
	if err != nil {
		logger.Errorf("failed to marshal log (%v):\n %v", err, gcpLog)
	}

	txn := instrumentation.StartTransaction("event_stackdriver_log_transaction", nil, nil)
	ogiproducer.Produce(GoogleLogTopic, logbytes, severity)
	instrumentation.EndTransaction(&txn)
}

func (gcp *GCPStackdriverLogsConsumer) Consume() {
	ctx := context.Background()

	adminClient, err := logadmin.NewClient(ctx,
		GoogleProjectId,
		option.WithCredentialsFile(GoogleCloudCredentialFile))
	if err != nil {
		logger.Fatalf("Failed to create logadmin client: %v", err)
	}

	iter := adminClient.Entries(ctx,
		logadmin.Filter(fmt.Sprintf(`logName = "projects/%s/logs/%s"`,
			GoogleProjectId, GoogleLogTopic)),
		logadmin.NewestFirst(),
	)

	for {
		entry, err := iter.Next()
		if err == iterator.Done {
			logger.Infoln("done")
		}
		if err != nil {
			logger.Fatalf("Failed: %v", err)
		}
		go gcp.produce(entry)
	}
}

func Consume() {
	gcp.Consume()
}
