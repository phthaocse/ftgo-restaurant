package message

import (
	"ftgo-restaurant/internal/outbound/interface/logger"
	"ftgo-restaurant/pkg/message"
	kafkaProducer "ftgo-restaurant/pkg/producer/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type RestaurantMessageProducer struct {
	Producer *kafkaProducer.Producer
	Logger   logger.Logger
}

func NewRestaurantMessageProducer(producer *kafkaProducer.Producer, logger logger.Logger) *RestaurantMessageProducer {
	return &RestaurantMessageProducer{
		Producer: producer,
		Logger:   logger,
	}
}

func (p *RestaurantMessageProducer) Produce(destination string, message message.Message) {
	err := p.Producer.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &destination, Partition: kafka.PartitionAny},
		Key:            message.Header,
		Value:          message.Payload,
	}, nil)
	if err != nil {
		return
	}
}
