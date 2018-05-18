package main

import (
	"os"

	ogiproducer "github.com/gojekfarm/ogi/producer"
)

func osPathExists(ospath string) (err error) {
	existsOrNot := "missing"
	if stat, err := os.Stat(ospath); err == nil {
		existsOrNot = "file"
		if stat.IsDir() {
			existsOrNot = "directory"
		}
	}
	ogiproducer.Produce(ospath, []byte(existsOrNot), "")
	return
}

func Transform(msg []byte) (err error) {
	return osPathExists(string(msg))
}
