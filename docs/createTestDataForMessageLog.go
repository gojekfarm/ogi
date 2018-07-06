package main

import (
	"fmt"

	logger "github.com/gojektech/ogi/logger"
	ogiproducer "github.com/gojektech/ogi/producer"
)

var (
	msgMap = map[string][]string{
		"nietzsche": []string{
			"without music life would be a mistake",
			"he who has a why to live, can bear almost any how",
			"the doer alone learneth",
		},
		"bob-dylan": []string{
			"no one is free, even the bords are chained in the sky",
			"a mistake is to commit a misunderstanding",
		},
		"william-s-burrough": []string{
			"nothing is true, everything is permitted",
			"paranoid is someone who knows a little of what's going on",
		},
	}
)

func main() {
	logger.SetupLogger()
	for msgTopic, messages := range msgMap {
		for idx, msg := range messages {
			msg = fmt.Sprintf("%s,%d,%s", msgTopic, idx, msg)
			ogiproducer.Produce("testdata", []byte(msg), msgTopic)
		}
	}
}
