package ogiproducer

import (
	"path"
	"plugin"

	"github.com/abhishekkr/gol/golenv"
	logger "github.com/gojekfarm/ogi/logger"
)

type ProducerPlugin struct {
	Name                   string
	NewProducerFunc        plugin.Symbol
	CloseFunc              plugin.Symbol
	GetMetadataFunc        plugin.Symbol
	GetPartitionNumberFunc plugin.Symbol
	ProduceMessageFunc     plugin.Symbol
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

	getMetadataFunc, err := p.Lookup("GetMetadata")
	if err != nil {
		logger.Fatalf("Error looking up 'GetMetadata': %s", err)
	}

	getPartitionNumberFunc, err := p.Lookup("GetPartitionNumber")
	if err != nil {
		logger.Fatalf("Error looking up 'GetPartitionNumber': %s", err)
	}

	producerMessageFunc, err := p.Lookup("ProduceMessage")
	if err != nil {
		logger.Fatalf("Error looking up 'ProduceMessage': %s", err)
	}

	return &ProducerPlugin{
		Name:                   path.Base(ProducerPluginPath),
		NewProducerFunc:        newProducerFunc,
		CloseFunc:              closeFunc,
		GetMetadataFunc:        getMetadataFunc,
		GetPartitionNumberFunc: getPartitionNumberFunc,
		ProduceMessageFunc:     producerMessageFunc,
	}
}

func (plugin *ProducerPlugin) NewProducer() {
	plugin.NewProducerFunc.(func())()
}

func (plugin *ProducerPlugin) Close() {
	plugin.CloseFunc.(func())()
}

func (plugin *ProducerPlugin) GetMetadata() {
	plugin.GetMetadataFunc.(func())()
}

func (plugin *ProducerPlugin) GetPartitionNumber(topic string, messageKey string) (partitionNumber int32) {
	return plugin.GetPartitionNumberFunc.(func(string, string) int32)(topic, messageKey)
}

func (plugin *ProducerPlugin) ProduceMessage(topic string, message []byte, partitionNumber int32) {
	plugin.ProduceMessageFunc.(func(string, []byte, int32))(topic, message, partitionNumber)
}
