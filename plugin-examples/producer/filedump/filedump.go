package main

import (
	"log"
	"os"
	"path"

	"github.com/abhishekkr/gol/golenv"
	logger "github.com/gojektech/ogi/logger"
)

type Filedump struct {
}

var (
	OgiFiledumpBasedir = golenv.OverrideIfEnv("OGI_FILEDUMP_BASEDIR", "/tmp")

	fd *Filedump
)

func init() {
	fd = new(Filedump)
}

func (filedump *Filedump) Close() {
	return
}

func (filedump *Filedump) Produce(topic string, message []byte, messageKey string) {
	filedir := path.Join(OgiFiledumpBasedir, topic)
	if err := os.MkdirAll(filedir, 0755); err != nil {
		logger.Errorf("skipping file creation, failed to create dir (%s):\n %v", filedir, err)
		return
	}

	filepath := path.Join(filedir, messageKey)
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(message); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	f.Sync()
}

func Close() {
	fd.Close()
}

func Produce(topic string, message []byte, messageKey string) {
	fd.Produce(topic, message, messageKey)
}
