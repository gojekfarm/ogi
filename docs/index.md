
## Ogi

utility to enable flexible [ETL](https://en.wikipedia.org/wiki/Extract,_transform,_load) scenarios

> ![ogi means a japanese fan](./ogi.png "ogi means a japanese fan")

---

#### How it works?

It contains 3 primary components. A consumer, transformer and producer.

Ogi initializes Consumer and let it handle the flow to transformer, or if required directly producer. That flow is internal to consumer used and not mandated. Similar internal flow freedom is granted to transformer and producer. Like, if required your producer can have multiple outputs anywhere from Kafka, S3 to something like e-mail.

All 3, `consumer`, `transformer` and `producer` are instantiated as per config and thus any combination of available types could be brought into play.

Consumer, Transformer and Producer support usage of `golang plugin`, so separately managed and developed constructs could be used in combination as well.
Since they are loaded as per configuration provided identifier, one can have combination of say built-in Kafka consumer with built-in kubernetes-log transformer but custom external plug-in of Google Cloud Datastore for cold storage of logs.
This also gives capability to write a producer sending output to more than one output sinks in same flow to achieve replication.

Can be scaled up easily using kubernetes/nomad/mesos/\* elastic deployments as it inherently has no context.

---

#### Concepts

* [design in detail](./design)

* [in-built workflows available](./types)

* [what is a golang plug-in and how can I write one for Ogi](./plugins)

---

#### Examples

* [a very simple usecase to understand design and writing plugins](./example-usecase-01.md)

---

#### Quickstart

* [get it running locally](./run-locally)

* [set of configurations to make ogi work to specific behavior](./config-set)

---

