package main

import (
	"math/rand"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	prototube "github.com/tiantiandas/prototube-go"
)

var (
	demoKafkaTopic               = "prototube_demo"
	demoKafkaBootstrapBrokerList = []string{"localhost:9092"}
)

func getDefaultKafkaProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll       // Wait for all in-sync replicas to commit to ack
	config.Producer.Compression = sarama.CompressionSnappy // Compress messages
	config.Producer.Flush.Frequency = 1 * time.Second      // Flush batches every 1 second
	config.Producer.Flush.Bytes = 1 * 1024 * 1024          // The best-effort number of bytes needed to trigger a flush.
	config.Producer.Flush.Messages = 1000                  // The best-effort number of messages needed to trigger a flush.
	// Retry in 5 min
	config.Producer.Retry.Max = 30                   // Retry 30 times for producer requests
	config.Producer.Retry.Backoff = 10 * time.Second // Retry interval 10 secs for producer requests
	return config
}

func generateRandomEvent() (*ExamplePrototubeMessage, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	uuidBytes, err := uuid.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return &ExamplePrototubeMessage{
		Int32Field:  int32(rand.Intn(10000)),
		Int64Field:  rand.Int63n(int64(10000) + 10000),
		DoubleField: rand.Float64(),
		StringField: uuid.String(),
		BytesField:  uuidBytes,
	}, nil
}

func main() {
	log.Infoln("Trying to start example producer!")
	producerConfig := &prototube.ProducerConfig{
		KafkaBootstrapBrokerList: demoKafkaBootstrapBrokerList,
		KafkaProducerConfig:      getDefaultKafkaProducerConfig(),
	}
	producer, err := prototube.NewWithConfig(demoKafkaTopic, producerConfig)
	if err != nil {
		log.Fatalf("Failed to initialize stream producer with topic [ %v ] and producer config [ %v ]: %v", demoKafkaTopic, producerConfig, err)
	}
	defer producer.Stop()
	for i := 0; i < 100; i++ {
		msg, err := generateRandomEvent()
		if err != nil {
			log.Fatalf("Failed to generate random event: %v", err)
		} else {
			producer.Emit(msg)
		}
	}
	log.Infof("Emit all messages to Kafka")
}
