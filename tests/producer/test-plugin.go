package main

type TestProducerPlugin struct {
}

var (
	p *TestProducerPlugin
)

func init() {
	p = new(TestProducerPlugin)
}

func (k *TestProducerPlugin) NewProducer() {
	return
}

func (k *TestProducerPlugin) Close() {
	return
}

func (k *TestProducerPlugin) Produce(topic string, message []byte, messageKey string) {
	return
}

func NewProducer() {
	p.NewProducer()
}

func Close() {
	p.Close()
}

func Produce(topic string, message []byte, messageKey string) {
	p.Produce(topic, message, messageKey)
}
