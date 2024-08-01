package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"route/internal/app/config"
)

type KafkaConsumer struct {
	brokers  []string
	Consumer sarama.Consumer
	topic    string
}

func NewKafkaConsumer(cfg config.KafkaConfig) (*KafkaConsumer, error) {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = false
	consumerConfig.Consumer.Offsets.AutoCommit.Enable = true
	consumerConfig.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second

	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(cfg.BrokerList, consumerConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		brokers:  cfg.BrokerList,
		Consumer: consumer,
		topic:    cfg.Topic,
	}, err
}

func (kc *KafkaConsumer) ReadMessages() error {
	partitionList, err := kc.Consumer.Partitions(kc.topic)
	if err != nil {
		return err
	}

	for _, partition := range partitionList {
		pc, err := kc.Consumer.ConsumePartition(kc.topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		defer pc.Close()

		for message := range pc.Messages() {
			fmt.Printf("Message received: %s\n", string(message.Value))
		}
	}

	return nil
}
