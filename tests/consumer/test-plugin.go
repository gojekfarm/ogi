package main

type TestConsumerPlugin struct {
}

var (
	p *TestConsumerPlugin
)

func init() {
	p = new(TestConsumerPlugin)
}

func (k *TestConsumerPlugin) Consume() {
	return
}

func Consume() {
	p.Consume()
}
