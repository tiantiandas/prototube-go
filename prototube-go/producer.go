package prototube

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
)

var (
	prototubeMessageHeader   = []byte{0x50, 0x42, 0x54, 0x42}
	kafkaBootstrapBrokerList = []string{"localhost:9092"}
)

// Producer a event log producer that pushes events into Kinesis stream.
type Producer struct {
	kafkaProducer sarama.AsyncProducer
	topic         string
}

// New create a new producer for the specified topic
func New(topic string) (*Producer, error) {
	producerConfig := &ProducerConfig{
		KafkaBootstrapBrokerList: kafkaBootstrapBrokerList,
		KafkaProducerConfig:      getDefaultKafkaProducerConfig(),
	}
	return NewWithConfig(topic, producerConfig)
}

// NewWithConfig create a new producer for the specified topic
func NewWithConfig(topic string, producerConfig *ProducerConfig) (*Producer, error) {
	if producerConfig == nil {
		producerConfig = &ProducerConfig{
			KafkaBootstrapBrokerList: kafkaBootstrapBrokerList,
			KafkaProducerConfig:      getDefaultKafkaProducerConfig(),
		}
	}
	if producerConfig.KafkaProducerConfig == nil {
		producerConfig.KafkaProducerConfig = getDefaultKafkaProducerConfig()
	}
	kafkaProducer, err := sarama.NewAsyncProducer(producerConfig.KafkaBootstrapBrokerList, producerConfig.KafkaProducerConfig)
	if err != nil {
		configJSON, _ := json.Marshal(producerConfig.KafkaProducerConfig)
		log.Fatalf("Failed to start Sarama producer with topic [%v] and producer config [%v]: [%v]", producerConfig.KafkaBootstrapBrokerList, configJSON, err)
	}
	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range kafkaProducer.Errors() {
			log.Errorf("Failed to write msg entry: %v", err)
		}
	}()

	return &Producer{
		kafkaProducer: kafkaProducer,
		topic:         topic,
	}, nil
}

func getDefaultKafkaProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll       // Wait for all in-sync replicas to commit to ack
	config.Producer.Compression = sarama.CompressionSnappy // Compress messages
	config.Producer.Flush.Frequency = 60 * time.Second     // Flush batches every 1 minute
	config.Producer.Flush.Bytes = 1 * 1024 * 1024          // The best-effort number of bytes needed to trigger a flush.
	config.Producer.Flush.Messages = 1000                  // The best-effort number of messages needed to trigger a flush.
	// Retry in 5 min
	config.Producer.Retry.Max = 30                   // Retry 30 times for producer requests
	config.Producer.Retry.Backoff = 10 * time.Second // Retry interval 10 secs for producer requests
	return config
}

// Stop stop the producer
func (p *Producer) Stop() error {
	if err := p.kafkaProducer.Close(); err != nil {
		log.Errorf("Failed to shut down kafka producer cleanly: %v", err)
		return err
	}
	return nil
}

// Emit a new log in PB format to the stream with the current timestamp
func (p *Producer) Emit(msg proto.Message) error {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	uuidBytes, err := uuid.MarshalBinary()
	if err != nil {
		return err
	}

	bytes, err := p.encode(time.Now().UTC().Unix(), uuidBytes, msg)
	if err != nil {
		return err
	}

	p.kafkaProducer.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(bytes),
	}
	return err
}

func (p *Producer) encode(ts int64, uuid []byte, m proto.Message) ([]byte, error) {
	h := &PrototubeMessageHeader{
		Version: 1,
		Ts:      ts,
		Uuid:    uuid,
	}

	buf := make([]byte, len(prototubeMessageHeader))
	copy(buf, prototubeMessageHeader)

	b := proto.NewBuffer(buf)
	if err := b.EncodeMessage(h); err != nil {
		return nil, err
	}
	if err := b.EncodeMessage(m); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
