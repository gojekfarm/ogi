package main

import (
	ogiproducer "github.com/gojekfarm/ogi/producer"
)

type TestTransformerLog struct {
}

var (
	p *TestTransformerLog
)

func (msgLog *TestTransformerLog) Transform(msg string, producer ogiproducer.Producer) (err error) {
	return
}

func Transform(msg string, producer ogiproducer.Producer) error {
	return p.Transform(msg, producer)
}
