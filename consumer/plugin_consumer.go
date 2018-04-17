package ogiconsumer

import (
	"path"
	"plugin"

	"github.com/abhishekkr/gol/golenv"
	logger "github.com/gojekfarm/ogi/logger"
)

type ConsumerPlugin struct {
	Name                string
	ConfigureFunc       plugin.Symbol
	NewConsumerFunc     plugin.Symbol
	SubscribeTopicsFunc plugin.Symbol
	EventHandlerFunc    plugin.Symbol
	CloseFunc           plugin.Symbol
}

var (
	ConsumerPluginPath = golenv.OverrideIfEnv("CONSUMER_PLUGIN_PATH", "./consumer.so")
)

func NewConsumerPlugin() Consumer {
	p, err := plugin.Open(ConsumerPluginPath)
	if err != nil {
		logger.Fatalf("Error reading plugin at %s: %s", ConsumerPluginPath, err)
	}

	configureFunc, err := p.Lookup("Configure")
	if err != nil {
		logger.Fatalf("Error looking up 'Configure': %s", err)
	}

	newConsumerFunc, err := p.Lookup("NewConsumer")
	if err != nil {
		logger.Fatalf("Error looking up 'NewConsumer': %s", err)
	}

	subscriberTopicsFunc, err := p.Lookup("SubscribeTopics")
	if err != nil {
		logger.Fatalf("Error looking up 'SubscribeTopics': %s", err)
	}

	eventHandlerFunc, err := p.Lookup("EventHandler")
	if err != nil {
		logger.Fatalf("Error looking up 'EventHandler': %s", err)
	}

	closeFunc, err := p.Lookup("Close")
	if err != nil {
		logger.Fatalf("Error looking up 'Close': %s", err)
	}

	return &ConsumerPlugin{
		Name:                path.Base(ConsumerPluginPath),
		ConfigureFunc:       configureFunc,
		NewConsumerFunc:     newConsumerFunc,
		SubscribeTopicsFunc: subscriberTopicsFunc,
		EventHandlerFunc:    eventHandlerFunc,
		CloseFunc:           closeFunc,
	}
}

func (plugin *ConsumerPlugin) Configure() {
	plugin.ConfigureFunc.(func())()
}

func (plugin *ConsumerPlugin) NewConsumer() {
	plugin.NewConsumerFunc.(func())()
}

func (plugin *ConsumerPlugin) SubscribeTopics() {
	plugin.SubscribeTopicsFunc.(func())()
}

func (plugin *ConsumerPlugin) EventHandler() {
	plugin.EventHandlerFunc.(func())()
}

func (plugin *ConsumerPlugin) Close() {
	plugin.CloseFunc.(func())()
}
