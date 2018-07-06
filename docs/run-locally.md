
### Run [Ogi](https://github.com/gojektech/ogi) Locally

* latest release: [v1.0](https://github.com/gojektech/ogi/releases/tag/v1.0)

```
## create ogi-conf.env with required configurations in $PWD
## help could be taken from lnik above
#
## if any plugin used, that '*.so' should be present in $PWD

wget -c https://github.com/gojektech/ogi/releases/download/v1.0/ogi-linux-amd64

docker run -it --env-file ogi-conf.env $PWD:/opt/ogi ubuntu:16.04 /opt/ogi/ogi-linux-amd64
```

* this uses [golang plugins](https://golang.org/pkg/plugin/) for extensibility, currently supported on linux, utilize docker to run if using something else

* [ogi-v1.0-linux-amd64](https://github.com/gojektech/ogi/releases/download/v1.0/ogi-linux-amd64) could run from a Docker on non-linux platform

---
