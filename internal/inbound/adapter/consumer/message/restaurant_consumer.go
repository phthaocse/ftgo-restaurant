package message

import (
	"ftgo-restaurant/internal/outbound/interface/logger"
	kafkaConsumer "ftgo-restaurant/pkg/consumer/kafka"
	"ftgo-restaurant/pkg/message"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type RestaurantMessageConsumer struct {
	Consumer   *kafkaConsumer.KafkaConsumer
	Logger     logger.Logger
	Dispatcher message.Dispatcher
}

func NewRestaurantConsumer(logger logger.Logger) *RestaurantMessageConsumer {
	return &RestaurantMessageConsumer{
		Consumer: kafkaConsumer.NewConsumer(logger),
		Logger:   logger,
	}
}

func (oc *RestaurantMessageConsumer) ProcessMessage(msg *kafka.Message) error {
	dispatchMsg := kafkaConsumer.ToMessage(msg)
	oc.Dispatcher.Dispatch(dispatchMsg)
	return nil
}

func (oc *RestaurantMessageConsumer) Subscribe(subscriberId string, channels map[string]struct{}, dispatcher message.Dispatcher) {
	err := oc.Consumer.SubscriptTopics(kafkaConsumer.MapToSlice(channels))
	oc.Dispatcher = dispatcher
	if err != nil {
		oc.Logger.Errorf("can't subscript topic order %v", err)
		panic(err)
	}
}

func (oc *RestaurantMessageConsumer) Start() {
	oc.Consumer.ListenAndProcess(oc.ProcessMessage)
}
