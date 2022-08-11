package kafka

import (
	"encoding/json"
	"ftgo-restaurant/internal/outbound/interface/logger"
	"ftgo-restaurant/pkg/helpers"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
)

type Producer struct {
	Producer       *kafka.Producer
	topics         []string
	Config         kafka.ConfigMap
	Logger         logger.Logger
	topicPartition []kafka.TopicPartition
}

func (c *Producer) ReadConfig() {
	viper.SetConfigName("kafka")
	viper.SetConfigType("yaml")
	projectPath := helpers.ProjectPath()
	viper.AddConfigPath(projectPath + "config")
	if err := viper.ReadInConfig(); err != nil {
		c.Logger.Error(err.Error())
	}
	c.Config = kafka.ConfigMap{}
	for key, val := range viper.GetStringMap("kafka-producer") {
		c.Config[key] = val
	}
	c.Logger.Info(c.Config)
}

func NewProducer(logger logger.Logger) *Producer {
	producer := Producer{}
	producer.Logger = logger
	producer.ReadConfig()
	var err error
	configByte, _ := json.Marshal(producer.Config)
	producer.Logger.Info("Create event_publisher with config: ", string(configByte))
	producer.Producer, err = kafka.NewProducer(&producer.Config)
	if err != nil {
		producer.Logger.Errorf("Create event_publisher failed: %v", err)
	}
	producer.Logger.Info("Created event_publisher ", producer.Producer.String())
	go func() {
		for e := range producer.Producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.Errorf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					logger.Infof("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()
	return &producer
}
