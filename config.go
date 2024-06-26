package prototube

import sarama "github.com/IBM/sarama"

// ProducerConfig configs to create a Sarama Kafka Producer.
type ProducerConfig struct {

	// Kafka broker list
	KafkaBootstrapBrokerList []string

	// Kafka producer config for sarama
	KafkaProducerConfig *sarama.Config
}
