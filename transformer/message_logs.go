package ogitransformer

import (
	"strings"

	ogiproducer "github.com/gojekfarm/ogi/producer"
)

type MessageLog struct {
	Key   string
	Topic string
}

func (msgLog *MessageLog) Transform(msg string, producer ogiproducer.Producer) (err error) {
	msgTokens := strings.Split(msg, ",")

	msgLog.Topic = msgTokens[0]
	msgLog.Key = msgTokens[1]

	ogiproducer.Produce(producer,
		msgLog.Topic,
		[]byte(msg),
		msgLog.Key)
	return
}

func NewMessageLog() LogTransformer {
	return &MessageLog{}
}
