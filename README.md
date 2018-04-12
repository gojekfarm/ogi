## Ogi

> ![ogi means a japanese fan](docs/ogi.png "ogi means a japanese fan")
>
> initially written to fan-out bulk topic `labels[app:appname]` tagged logs pushed from Kubernetes to Kafka, into `app` specific topics
>
> evolved to be usable for flexible [ETL](https://en.wikipedia.org/wiki/Extract,_transform,_load) scenarios that can scaled up as multiple instances

---

#### way it works

* has a consumer, runs Subscribe on it... then calls its EventHandler
> consumer's EventHandler could choose desired flow for every event received

* events could pass details and desired producer for it to Transformer, transformer follows it's own computation and pass on updated state to producer
> if no transformation is required for this `consumer`, the Producer could be triggered directly as well

* producer just does that, produce message to required target

> here all 3, `consumer`, `transformer` and `producer` are instantiated as per config and thus any combination of available types could be brought into play

---

#### current scenarios available

* pull data from a kafka topic and write to multiple topics based on the hashing logic

---

## sub-packages

### consumer

> this package provides a `Consume(consumer)` method to be able to munch on any kafka-topic,
>
> currently `confluent-kafka` implementation for consumer is available and can be configured and passed

> topic subscriber currently handles
> * for each messages recieved, invoke `ogi's tansformer` with producer and value from message
> * assign partition
> * unassign partition
> * exit on error event


### producer

> this package provides a `Produce(producer, topic, message)` method to be able to churn out required message to any topic on provided producer,
>
> currently `confluent-kafka` implementation for producer is available which can be configured and passed

> `confluent-kafka` producer let's figure out a partition for message using CRC32 over total partition count, if can't uses PartitionAny conf
>
> then message and topic gets produced to the calculated partition using `ProduceMessage`


### transformer

> this package provides `Transform(producer, message)` method which calls delegates transform to configured `LogTransformer` for that process,
>
> currently `KubernetesKafkaLog` is available implementation of `LogTransformer`

> `KubernetesKafkaLog` which checks for Kubernetes.Labels for configured label to be picked as destination kafka topic
>
> it applies Kubernetes.PodName as message-key to be used and then produces message to passed through producer

---
