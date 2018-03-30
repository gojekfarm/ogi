## Kafka-Ogi

> ![ogi means a japanese fan](docs/ogi.png "ogi means a japanese fan")

service used to pull data from a kafka topic and write to multiple topics based on the hashing logic

---

## sub-packages

### consumer

> this package prvides a `Consume(consumer)` method to be able to munch on any kafka-topic,
>
> currently `confluent-kafka` implementation for consumer is available and can be configured and passed

> topic subscriber currently handles
> * for each messages recieved, invoke `ogi's tansformer` with producer and value from message
> * assign partition
> * unassign partition
> * exit on error event

---
