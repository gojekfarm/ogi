package main

type TestConsumerPlugin struct {
}

var (
	p *TestConsumerPlugin
)

func init() {
	p = new(TestConsumerPlugin)
}

func (k *TestConsumerPlugin) Configure() {
	return
}

func (k *TestConsumerPlugin) NewConsumer() {
	return
}

func (k *TestConsumerPlugin) SubscribeTopics() {
	return
}

func (k *TestConsumerPlugin) EventHandler() {
	return
}

func (k *TestConsumerPlugin) Close() {
	return
}

func Configure() {
	p.Configure()
}

func NewConsumer() {
	p.NewConsumer()
}

func SubscribeTopics() {
	p.SubscribeTopics()
}

func EventHandler() {
	p.EventHandler()
}

func Close() {
	p.Close()
}
