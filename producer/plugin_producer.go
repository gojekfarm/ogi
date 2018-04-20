package ogiproducer

import (
	"path"
	"plugin"

	"github.com/abhishekkr/gol/golenv"
	logger "github.com/gojekfarm/ogi/logger"
)

type ProducerPlugin struct {
	Name            string
	NewProducerFunc plugin.Symbol
	CloseFunc       plugin.Symbol
	ProduceFunc     plugin.Symbol
}

var (
	ProducerPluginPath = golenv.OverrideIfEnv("PRODUCER_PLUGIN_PATH", "./producer.so")
)

func NewProducerPlugin() Producer {
	p, err := plugin.Open(ProducerPluginPath)
	if err != nil {
		logger.Fatalf("Error reading plugin: %s", err)
	}

	newProducerFunc, err := p.Lookup("NewProducer")
	if err != nil {
		logger.Fatalf("Error looking up 'NewProducer': %s", err)
	}

	closeFunc, err := p.Lookup("Close")
	if err != nil {
		logger.Fatalf("Error looking up 'Close': %s", err)
	}

	produceFunc, err := p.Lookup("Produce")
	if err != nil {
		logger.Fatalf("Error looking up 'Produce': %s", err)
	}

	return &ProducerPlugin{
		Name:            path.Base(ProducerPluginPath),
		NewProducerFunc: newProducerFunc,
		CloseFunc:       closeFunc,
		ProduceFunc:     produceFunc,
	}
}

func (plugin *ProducerPlugin) NewProducer() {
	plugin.NewProducerFunc.(func())()
}

func (plugin *ProducerPlugin) Close() {
	plugin.CloseFunc.(func())()
}

func (plugin *ProducerPlugin) Produce(topic string, message []byte, messageKey string) {
	plugin.ProduceFunc.(func(string, []byte, string))(topic, message, messageKey)
}
