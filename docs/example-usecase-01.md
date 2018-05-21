
## First Example

### UseCase: Identify filesystem paths to be a file, dir or missing from a given list in file

This example is not an actual show of it's capability, but design.

Also given requirement of usecase, required plug-ins code could be analyzed much easier.

#### Transposing it to an ETL like solution

Extract the filesystem paths from each line.
Transform the result to check if it's a file, dir or missing.
Load the result to required output sink.

So we need a `consumer` plugin here to read a file and run `transformer` on it line by line.
The `transformer` plugin to do the checks and formulate required output to be passed to `producer`.
The `producer` just printing out the results on stdout for us to know.

Here we'll be using custom plugins for consumer, transformer and producer.

We'll have a quick look at their code to identify how a minimal plugin could be written as well.

> Ogi's primary concern in just invoking Consumer and then let it take the charge of flow of logic.
>
> Ogi doesn't require or dictate anything else. Then while using Consumer (or Transformer, Producer), if using a plugin they should also abide by simple one or two primary exported functions. Anything else in them is not of Ogi's concern or purview.

* `file_line_by_line` Consumer

> The only exported function required by consumer plugin is `Consume()`. There is no close export required as `Consume()` just gets called once and can manage it internally.
>
> As we'll notice in this case as desired it opens a file, read it line by line and passes every line individually to a different `Transform([]byte) error`.
>
> But as we'll notice here, `Consume()` doesn't need to call `Transform([]byte) err`, it could something else entirely if it wants. It could directly call producer, a separate external library method or just do everything by itself if it desires.
>
> We've also made plugin responsible to avail config for its own requirement. Although that need to be configured.

```
package main

import (
	"bufio"
	"log"
	"os"

	"github.com/abhishekkr/gol/golenv"

	ogitransformer "github.com/gojekfarm/ogi/transformer"
)

var FileToConsume = golenv.OverrideIfEnv("OGI_FILE_TO_CONSUME", "/tmp/ogi-consumed")

func Consume() {
	fileHandle, err := os.Open(FileToConsume)
	if err != nil { log.Fatalln(err) }
	defer fileHandle.Close()

	fileScanner := bufio.NewScanner(fileHandle)
	for fileScanner.Scan() {
		err = ogitransformer.Transform([]byte(fileScanner.Text()))
    if err != nil { log.Println("failed for:", err) }
	}
}
```

* `os_path_exists` Tranformer

> For a plugin of transformer Ogi expects a `Transform([] byte) error` exported function.
>
> Here it checks what message (type of filesystem path) it has received, passes details on to `Produce(string, []byte, string)`.
>
> Again, like consumer, Ogi doesn't mandate transformer to use producer. Transformer can choose a different flow it feels like.

```
package main

import (
	"os"

	ogiproducer "github.com/gojekfarm/ogi/producer"
)

func Transform(msg []byte) (err error) {
  ospath := string(msg)
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
```

* `echo` Producer

> To be a plugin for producer, it need to export two functions.
>
> * `Produce(string, []byte, string)` to call Produce on required `topic`, `message` and `message-key` in respective order.
>
> * `Close()` is exported as well to allow certain flows to keep an open connection right from consumer for entire lifetime of consumption and then call a final close at the end.
>
> We'll notice in given example only topic and message is passed by transformer, but an empty message-key. As there is no need of it.
> Now it is job of every producer to mandate what fields are mandatory for it to be provided and what can be empty.

```
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
```

---

#### Running above solution with Ogi

* requires compiling these plugins separately as

```
	cd plugin-examples/consumer/file_line_by_line; \
		go build -o "../../../out/consumer-file-line-by-line.so" -buildmode=plugin . ; \
		cd - ; echo "compiled consumer.file_line_by_line plugin"

	cd plugin-examples/transformer/os_path_exists ; \
		go build -o "../../../out/transformer-os-path-exists.so" -buildmode=plugin . ; \
		cd - ; echo "compiled transformer.os_path_exists plugin"

	cd plugin-examples/producer/echo; \
		go build -o "../../../out/producer-echo.so" -buildmode=plugin . ; \
		cd - ; echo "compiled producer.echo plugin"
```

* creating config so required plugins are loaded, say as `env.cfg`

```
export CONSUMER_TYPE="plugin" # type given if want a plugin to be used
export CONSUMER_PLUGIN_PATH=$(pwd)"/consumer-file-line-by-line.so"  # path to plugin required if consumer type is plugin
## required by custom consumer plugin file-line-by-line config
export OGI_FILE_TO_CONSUME="demo01.data"

export PRODUCER_TYPE="plugin"
export PRODUCER_PLUGIN_PATH=$(pwd)"/producer-echo.so"

export TRANSFORMER_TYPE="plugin"
export TRANSFORMER_PLUGIN_PATH=$(pwd)"/transformer-os-path-exists.so"

export NEWRELIC_APP_NAME="ogi-test"
export NEWRELIC_LICENSE_KEY="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

export LOG_LEVEL="error"
```

* providing whatever plug-in internally desires (in this case consumer needs a file path list present), say as `demo01.data`

```
/tmp
/tmp/this-show-be-file
/tmp/missing
/bin
/bin/sh
```

* then just running Ogi, should be similar to

```
$ source ./env.cfg ; ./ogi

topic: /tmp
directory
topic: /tmp/this-show-be-file
file
topic: /tmp/missing
missing
topic: /bin
directory
topic: /bin/sh
file
```

---
