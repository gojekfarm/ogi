package main

var (
	p *TestProducerPlugin
)

type TestProducerPlugin struct {
}

func (k *TestProducerPlugin) NewProducer() {
	return
}

func (k *TestProducerPlugin) Close() {
	return
}

func (k *TestProducerPlugin) GetMetadata() {
	return
}

func (k *TestProducerPlugin) GetPartitionNumber(topic string, messageKey string) (partitionNumber int32) {
	return
}
func (k *TestProducerPlugin) ProduceMessage(topic string, message []byte, partitionNumber int32) {
	return
}

func NewProducer() {
	p.NewProducer()
}

func Close() {
	p.Close()
}

func GetMetadata() {
	p.GetMetadata()
}

func GetPartitionNumber(topic string, messageKey string) (partitionNumber int32) {
	return p.GetPartitionNumber(topic, messageKey)
}

func ProduceMessage(topic string, message []byte, partitionNumber int32) {
	p.ProduceMessage(topic, message, partitionNumber)
}
