package kafka

import (
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"route/internal/app/config"
)

type KafkaProducer struct {
	brokers  []string
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(cfg config.KafkaConfig) (*KafkaProducer, error) {
	producerConfig := sarama.NewConfig()

	producerConfig.Producer.Partitioner = sarama.NewRandomPartitioner

	producerConfig.Producer.RequiredAcks = sarama.WaitForLocal

	producerConfig.Producer.Idempotent = true
	producerConfig.Net.MaxOpenRequests = 1

	producerConfig.Producer.CompressionLevel = sarama.CompressionLevelDefault

	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.Return.Errors = true

	producerConfig.Producer.Compression = sarama.CompressionGZIP

	syncProducer, err := sarama.NewSyncProducer(cfg.BrokerList, producerConfig)
	if err != nil {
		return &KafkaProducer{}, errors.Wrap(err, "error with sync kafka-producer")
	}

	return &KafkaProducer{
		brokers:  cfg.BrokerList,
		producer: syncProducer,
		topic:    cfg.Topic,
	}, nil
}

func (kp *KafkaProducer) SendMessage(message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: kp.topic,
		Value: sarama.ByteEncoder(message),
	}
	_, _, err := kp.producer.SendMessage(msg)
	return err
}
