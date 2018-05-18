package main

import (
	"fmt"
	"log"

	"github.com/abhishekkr/gol/golenv"
)

var (
	Separator = golenv.OverrideIfEnv("OGI_ECHO_SEPARATOR", "")
)

func Close() {
	fmt.Println("~ ogi is done printing data")
}

func Produce(topic string, message []byte, messageKey string) {
	if topic != "" {
		fmt.Println("topic:", topic)
	}
	if messageKey != "" {
		fmt.Println("key:", messageKey)
	}
	if len(message) != 0 {
		fmt.Println(string(message))
	} else {
		log.Println("received blank message")
	}
}
