
### Dev Env Help


* git clone `kafka-ogi` repo

```
git clone https://github.com/gojekfarm/kafka-ogi
```


* fetch dependency manager

```
go get -u github.com/golang/dep/cmd/dep
```

* fetch dependencies

```
dep ensure
```

* prepare environment config file

```
cp env.sample env
## now replace values for all keys there with required one
```

* run tests

```
dep ensure
source env.sample
go test -gcflags=-l github.com/gojekfarm/kafka-ogi/consumer
```

---
