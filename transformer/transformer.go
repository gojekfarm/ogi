package ogitransformer

import (
	"fmt"

	"github.com/abhishekkr/gol/golenv"

	"github.com/gojektech/ogi/instrumentation"

	logger "github.com/gojektech/ogi/logger"
)

type Transformer interface {
	Transform([]byte) error
}

type NewTransformer func() Transformer

var (
	TransformerType = golenv.OverrideIfEnv("TRANSFORMER_TYPE", "kubernetes-kafka-log")

	transformerMap = map[string]NewTransformer{
		"kubernetes-kafka-log": NewKubernetesKafkaLog,
		"plugin":               NewTransformerPlugin,
	}
)

func init() {
	validateConfig()
}

func validateConfig() {
	var missingVariables string
	if KubernetesTopicLabel == "" {
		logger.Warn("Missing Env Config: 'PRODUCER_KUBERNETES_TOPIC_LABEL', can't use Kubernetes Label based transformers")
	}

	if TransformerType == "" {
		missingVariables = fmt.Sprintf("%s TRANSFORMER_TYPE", missingVariables)
	}

	if missingVariables != "" {
		logger.Fatalf("Missing Env Config:%s", missingVariables)
	}
}

func Transform(msg []byte) {
	txn := instrumentation.StartTransaction("transform_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	transformer := transformerMap[TransformerType]()
	if err := transformer.Transform(msg); err != nil {
		// produce to dead-man-talking topic
		logger.Warn(err)
	}
}
