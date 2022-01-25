Prototube
=========

Package prototube implements an API to publish strongly-typed events into Kafka.

Applications can publish events into a specific topic. A topic is always associated with a schema
which defines the schema of the event. Events that do not conform with the schemas are rejected.

Internally prototube encodes the events using the following format:

```
<Magic Number> <Header> <Event>
```

* Magic Number: 0x50, 0x42, 0x54, 0x42
* Header: Protobuf-encoded structure that contains the metadata of the event (e.g, timestamp / uuid)
* Event: Protobuf-encoded event

How to use
==========

Prerequisite
---------

Please follow this [Kafka Quickstart](https://kafka.apache.org/quickstart) link to install and start Kafka locally.

Quickstart
---------

Build and run the example application to produce random messages to local Kafka.

```sh
$ cd examples
$ go run main.go
```

Compile a proto file
--------------------

Below command generates the example `example.pb.go` from `example.proto` with module `main`.

```sh
$ cd examples
$ protoc --go_out=. idl/example.proto
```

Use prototube producer
----------------------

Please see this [example](https://github.com/tiantiandas/prototube-go/blob/master/examples/main.go) for your reference.

Code snippet:
```go
producer, err := prototube.NewWithConfig("testTopic", &prototube.ProducerConfig{
	KafkaBootstrapBrokerList: []string{"localhost:9092"},
})

producer.Emit(&ExamplePrototubeMessage{
	Int32Field:  int32(rand.Intn(10000)),
	Int64Field:  rand.Int63n(int64(10000) + 10000),
	DoubleField: rand.Float64(),
})
```
