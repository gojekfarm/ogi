
## config set

* main set of configs specifying which combination of consumer, transformer and producer are used

```
CONSUMER_TYPE="confluent-kafka"
CONSUMER_PLUGIN_PATH="${ABSOLUTE_PATH_FOR_CONSUMER_PLUGIN}" ## required if CONSUMER_TYPE="plugin"

PRODUCER_TYPE="confluent-kafka"
PRODUCER_PLUGIN_PATH="${ABSOLUTE_PATH_FOR_PRODUCER_PLUGIN}" ## required if PRODUCER_TYPE="plugin"

TRANSFORMER_TYPE="message-log"
TRANSFORMER_PLUGIN_PATH="${ABSOLUTE_PATH_FOR_TRANSFORMER_PLUGIN}" ## required if TRANSFORMER_TYPE="plugin"
```


* if using consumer 'confluent-kafka'

```
CONSUMER_KAFKA_TOPICS="testdata"
CONSUMER_BOOTSTRAP_SERVERS="127.0.0.1:9092"
CONSUMER_GROUP_ID="ogi-id"
CONSUMER_SESSION_TIMEOUT_MS="6000"
CONSUMER_GOEVENTS_CHANNEL_ENABLE="true"
CONSUMER_GOEVENTS_CHANNEL_SIZE="1000"
CONSUMER_GO_APPLICATION_REBALANCE_ENABLE="true"
```


* if using producer 'confluent-kafka'

```
PRODUCER_BOOTSTRAP_SERVERS="127.0.0.1:9092"
PRODUCER_KAFKA_TOPIC_LABEL="app"
```


* instrumentation and logging

```
NEWRELIC_APP_NAME="ogi-test"
NEWRELIC_LICENSE_KEY="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

LOG_LEVEL="info"
```

---
