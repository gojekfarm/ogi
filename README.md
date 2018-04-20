## Ogi

> ![ogi means a japanese fan](docs/ogi.png "ogi means a japanese fan")
>
> initially written to fan-out bulk topic `labels[app:appname]` tagged logs pushed from Kubernetes to Kafka, into `app` specific topics
>
> evolved to be usable for flexible [ETL](https://en.wikipedia.org/wiki/Extract,_transform,_load) scenarios that can scaled up as multiple instances

---

> _this uses [golang plugins](https://golang.org/pkg/plugin/) for extensibility, currently supported on linux, utilize docker to run if using something else_
> [ogi-v1.0-linux-amd64](https://github.com/gojekfarm/ogi/releases/download/v1.0/ogi-linux-amd64) could run from a Docker on non-linux platform

* latest release: [v1.0](https://github.com/gojekfarm/ogi/releases/tag/v1.0)

[set of configurations to make ogi work to specific behavior](./docs/config-set.md)

---

* [design in detail](./docs/design.md)

* [in-built workflows available](./docs/types.md)

---

#### Quickstart

```
## create ogi-conf.env with required configurations in $PWD
## help could be taken from lnik above
#
## if any plugin used, that '*.so' should be present in $PWD

wget -c https://github.com/gojekfarm/ogi/releases/download/v1.0/ogi-linux-amd64

docker run -it --env-file ogi-conf.env $PWD:/opt/ogi ubuntu:16.04 /opt/ogi/ogi-linux-amd64
```

---
