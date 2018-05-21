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
	err := ogitransformer.Transform([]byte(lyne))
	if err != nil {
		log.Println("failed for:", err)
	}
}

func consumeFile() {
	fileHandle, err := os.Open(FileToConsume)
	if err != nil {
		log.Fatalln(err)
	}
	defer fileHandle.Close()

	fileScanner := bufio.NewScanner(fileHandle)
	for fileScanner.Scan() {
		transform(fileScanner.Text())
	}
}

func Consume() {
	consumeFile()
}
