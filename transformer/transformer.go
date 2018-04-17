package ogitransformer

import (
	"fmt"

	"github.com/abhishekkr/gol/golenv"

	"github.com/gojekfarm/ogi/instrumentation"

	logger "github.com/gojekfarm/ogi/logger"
	ogiproducer "github.com/gojekfarm/ogi/producer"
)

type Transformer interface {
	Transform(string, ogiproducer.Producer) error
}

type NewTransformer func() Transformer

var (
	TransformerType = golenv.OverrideIfEnv("TRANSFORMER_TYPE", "kubernetes-kafka-log")

	transformerMap = map[string]NewTransformer{
		"kubernetes-kafka-log": NewKubernetesKafkaLog,
		"plugin":               NewTransformerPlugin,
	}
)

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

func Transform(producer ogiproducer.Producer, msg string) {
	txn := instrumentation.StartTransaction("transform_transaction", nil, nil)
	defer instrumentation.EndTransaction(&txn)

	transformer := transformerMap[TransformerType]()
	if err := transformer.Transform(msg, producer); err != nil {
		// produce to dead-man-talking topic
		logger.Warn(err)
	}
}
