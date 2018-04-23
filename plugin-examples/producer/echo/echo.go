package main

import "fmt"

type Echo struct {
}

var (
	e *Echo
)

func init() {
	e = new(Echo)
}

func (echo *Echo) Close() {
	fmt.Println("~ ogi is done printing data")
}

func (echo *Echo) Produce(topic string, message []byte, messageKey string) {
	fmt.Println("topic:", topic, "; key:", messageKey)
	fmt.Println(string(message))
	fmt.Println("*******************************************************")
}

func Close() {
	e.Close()
}

func Produce(topic string, message []byte, messageKey string) {
	e.Produce(topic, message, messageKey)
}
