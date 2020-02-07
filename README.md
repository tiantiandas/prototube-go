Prototube
=========

Package prototube implements an API to publish strongly-typed events into Kafka.

Applications can publish events into a specific topic. A topic is always associated with a schema
which defines the schema of the event. Events that do not conform with the schemas are rejected.

Internally prototube encodes the events using the following format:

> \<Magic Number\> \<Header\> \<Event\>

* Magic Number: 0x50, 0x42, 0x54, 0x42
* Header: Protobuf-encoded structure that contains the metadata of the event (e.g, timestamp / uuid)
* Event: Protobuf-encoded event

How to use
==========

Prequiste
---------

Please follow this [Kafka Quickstart](https://kafka.apache.org/quickstart) link to install and start Kafka locally.

Quickstart
---------

Build and run the example application to produce random messages to local Kafka.

```
cd examples
go build .
./example
```

Compile a proto file
--------------------

Below command generates the example `example.pb.go` from `example.proto` with module `main`.

```
cd examples
protoc --go_out=import_path=main:. -Iidl idl/example.proto
```

Use prototube producer
----------------------

Please see this [example](https://github.com/fx19880617/prototube-go/blob/master/examples/main.go) for your reference.

Code snippet:
```
producer, err := prototube.NewWithConfig("testTopic", &prototube.ProducerConfig{
	KafkaBootstrapBrokerList: []string{"localhost:9092"},
})

producer.Emit(&ExamplePrototubeMessage{
	Int32Field:  int32(rand.Intn(10000)),
	Int64Field:  rand.Int63n(int64(10000) + 10000),
	DoubleField: rand.Float64(),
})
```