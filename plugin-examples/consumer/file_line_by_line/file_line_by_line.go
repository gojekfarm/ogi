package main

import (
	"bufio"
	"log"
	"os"

	"github.com/abhishekkr/gol/golenv"

	ogitransformer "github.com/gojekfarm/ogi/transformer"
)

var (
	FileToConsume = golenv.OverrideIfEnv("OGI_FILE_TO_CONSUME", "/tmp/ogi-consumed")
)

func transform(lyne string) {
	ogitransformer.Transform([]byte(lyne))
}

func consumeFile() {
	fileHandle, err := os.Open(FileToConsume)
	if err != nil {
		log.Fatalln(err)
	}
	defer fileHandle.Close()

	fileScanner := bufio.NewScanner(fileHandle)
	for fileScanner.Scan() {
		ogitransformer.Transform([]byte(fileScanner.Text()))
	}
}

func Consume() {
	consumeFile()
}
