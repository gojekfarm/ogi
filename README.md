## Ogi

[![Go Report Card](https://goreportcard.com/badge/gojektech/ogi)](https://goreportcard.com/report/gojektech/ogi)


utility to enable flexible [ETL](https://en.wikipedia.org/wiki/Extract,_transform,_load) scenarios

> ![ogi means a japanese fan](docs/ogi.png "ogi means a japanese fan")
>
> Initially written to fan-out bulk topic `labels[app:appname]` tagged logs pushed from Kubernetes to Kafka, into `app` specific topics.
>
> Can be scaled up using Kubernetes/Nomad/Mesos elastic deployments, it inherently has no context.


---

#### How it works?

It contains 3 primary components. A consumer, transformer and producer.

Ogi initializes Consumer and let it handle the flow to transformer, or if required directly producer. That flow is internal to consumer used and not mandated. Similar internal flow freedom is granted to transformer and producer. Like, if required your producer can have multiple outputs anywhere from Kafka, S3 to something like e-mail.

All 3, `consumer`, `transformer` and `producer` are instantiated as per config and thus any combination of available types could be brought into play.

Consumer, Transformer and Producer support usage of `golang plugin`, so separately managed and developed constructs could be used in combination as well.
Since they are loaded as per configuration provided identifier, one can have combination of say built-in Kafka consumer with built-in kubernetes-log transformer but custom external plug-in of Google Cloud Datastore for cold storage of logs.
This also gives capability to write a producer sending output to more than one output sinks in same flow to achieve replication.

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

wget -c https://github.com/gojektech/ogi/releases/download/v1.0/ogi-linux-amd64

docker run -it --env-file ogi-conf.env $PWD:/opt/ogi ubuntu:16.04 /opt/ogi/ogi-linux-amd64
```

[set of configurations to make ogi work to specific behavior](./docs/config-set.md)

> _this uses [golang plugins](https://golang.org/pkg/plugin/) for extensibility, currently supported on linux, utilize docker to run if using something else_
> [ogi-v1.0-linux-amd64](https://github.com/gojektech/ogi/releases/download/v1.0/ogi-linux-amd64) could run from a Docker on non-linux platform

* latest release: [v1.0](https://github.com/gojektech/ogi/releases/tag/v1.0)

---
