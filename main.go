package main

import (
	"fmt"

	ogiconsumer "github.com/gojekfarm/ogi/consumer"
	logger "github.com/gojekfarm/ogi/logger"
)

func main() {
	logger.SetupLogger()
	fmt.Println(`
	   oooo          ggg      iiii
	 ooo  ooo      gg  gg       ii
	000    000       gggg       ii
	 ooo  ooo           gg      ii
	   oooo         gggg        ii
	`)
	ogiconsumer.Consume()
}
