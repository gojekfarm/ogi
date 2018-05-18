package main

import (
	"strings"

	logger "github.com/gojekfarm/ogi/logger"
	ogiproducer "github.com/gojekfarm/ogi/producer"
)

type MessageLog struct {
	Key   string
	Topic string
}

var (
	messageLog *MessageLog
)

func init() {
	messageLog = new(MessageLog)
}

func (msgLog *MessageLog) Transform(msg string) (err error) {
	logger.Infoln("message recieved is", msg)
	msgTokens := strings.Split(msg, ",")

	if len(msgTokens) < 3 {
		logger.Warnf("skipping msg due to invalid format : %s", msg)
		return
	}

	logger.Infoln(msgTokens, len(msgTokens), msgLog)
	msgLog.Topic = msgTokens[0]
	msgLog.Key = msgTokens[1]

	ogiproducer.Produce(msgLog.Topic,
		[]byte(msg),
		msgLog.Key)
	return
}

func Transform(msg []byte) (err error) {
	return messageLog.Transform(string(msg))
}
