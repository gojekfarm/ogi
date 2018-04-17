
### Dev Env Help


* git clone `kafka-ogi` repo

```
git clone https://github.com/gojekfarm/kafka-ogi
```


* fetch dependency manager

```
go get -u github.com/golang/dep/cmd/dep

## or

make setup
```


* fetch dependencies

```
dep ensure

## or

make build-deps
```


* prepare environment config file

```
cp env.sample env
## now replace values for all keys there with required one
```


* run tests

```
make build-deps ; source tests/tests-env.cfg ; make test
```

---
