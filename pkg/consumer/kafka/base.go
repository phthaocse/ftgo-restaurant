package kafka

import (
	"encoding/json"
	"fmt"
	"ftgo-restaurant/internal/outbound/interface/logger"
	"ftgo-restaurant/pkg/helpers"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const MaxPartition = 50

type KafkaConsumer struct {
	Consumer       *kafka.Consumer
	topics         []string
	Config         kafka.ConfigMap
	Logger         logger.Logger
	topicPartition []kafka.TopicPartition
}

func NewConsumer(logger logger.Logger) *KafkaConsumer {
	consumer := KafkaConsumer{}
	consumer.Logger = logger
	consumer.ReadConfig()
	var err error
	configByte, _ := json.Marshal(consumer.Config)
	consumer.Logger.Info("Create consumer with config: ", string(configByte))
	consumer.Consumer, err = kafka.NewConsumer(&consumer.Config)
	if err != nil {
		consumer.Logger.Errorf("Create consumer failed: %v", err)
	}
	consumer.Logger.Info("Created consumer ", consumer.Consumer.String())
	return &consumer
}

func (c *KafkaConsumer) ReadConfig() {
	viper.SetConfigName("kafka")
	viper.SetConfigType("yaml")
	projectPath := helpers.ProjectPath()
	viper.AddConfigPath(projectPath + "config")
	if err := viper.ReadInConfig(); err != nil {
		c.Logger.Error(err.Error())
	}
	c.Config = kafka.ConfigMap{}
	for key, val := range viper.GetStringMap("kafka-consumer") {
		c.Config[key] = val
	}
	c.Logger.Info(c.Config)
}

func (c *KafkaConsumer) rebalanceCallback(consumer *kafka.Consumer, event kafka.Event) error {
	switch event.(type) {
	case kafka.AssignedPartitions:
		assignedPartitions := event.(kafka.AssignedPartitions)
		c.Logger.Infof("Partitions were assigned: %v", assignedPartitions.Partitions)
		c.topicPartition = assignedPartitions.Partitions
	case kafka.RevokedPartitions:
		revokedPartitions := event.(kafka.RevokedPartitions)
		c.Logger.Infof("Partitions were revoked: %v", revokedPartitions.Partitions)
		topicPartition, err := consumer.Commit()
		if err != nil {
			c.Logger.Errorf("commit failed %v", err)
		}
		c.topicPartition = topicPartition
	}
	return nil
}

func (c *KafkaConsumer) SubscriptTopic(topic string) error {
	c.topics = append(c.topics, topic)
	return c.Consumer.SubscribeTopics(c.topics, c.rebalanceCallback)
}

func (c *KafkaConsumer) SubscriptTopics(topics []string) error {
	c.topics = append(c.topics, topics...)
	return c.Consumer.SubscribeTopics(c.topics, c.rebalanceCallback)
}

func (c *KafkaConsumer) SubscriptAllTopics() error {
	return c.Consumer.SubscribeTopics(c.topics, c.rebalanceCallback)
}

type ProcessMessageFn func(message *kafka.Message) error

func (c *KafkaConsumer) processMessage(messageChans []chan *kafka.Message, processFn ProcessMessageFn) {
	for _, messageChan := range messageChans {
		go func(messageChan <-chan *kafka.Message) {
			for msg := range messageChan {
				c.Logger.Infof("start processing message: %v, key: %s, value: %s, topicpartion: %v", msg, string(msg.Key), string(msg.Value), msg.TopicPartition)
				retry := 0
				for retry < 10 {
					retry++
					if err := processFn(msg); err != nil {
						c.Logger.Errorf("process message error: %v, key: %s, value: %s, topicpartion: %v with error %v", msg, string(msg.Key), string(msg.Value), msg.TopicPartition, err)
						time.Sleep(time.Second)
					}
					break
				}
				go func(message *kafka.Message) {
					_, err := c.Consumer.StoreMessage(message)
					if err != nil {
						c.Logger.Errorf("commit failed due to %v, the next commit will serve as retry")
						return
					}
					c.Logger.Infof("store message %v successfully", message)
				}(msg)
			}
		}(messageChan)
	}
}

func (c *KafkaConsumer) getCurrentPartition() (int, error) {
	numPartition := 0
	var consumerMetaData *kafka.Metadata
	var err error

	numRetry := 0
	for numRetry < 10 {
		numRetry++
		if numRetry >= 10 {
			return 0, err
		}
		consumerMetaData, err = c.Consumer.GetMetadata(nil, true, 10000)
		if err != nil {
			c.Logger.Errorf("Can't get consumer data due to error %v", err)
		}
		break
	}

	for _, topicMetaData := range consumerMetaData.Topics {
		if topicMetaData.Topic != "__consumer_offsets" {
			numPartition += len(topicMetaData.Partitions)
		}
	}
	return numPartition, nil
}

func (c *KafkaConsumer) ListenAndProcess(processFn ProcessMessageFn) {
	defer c.Consumer.Close()
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	var numChan int
	numPartition, err := c.getCurrentPartition()
	if err != nil {
		c.Logger.Infof("Due to the metadata couldn't be gotten from broker, the default max partition will be used")
		numChan = MaxPartition
	} else {
		numChan = numPartition
	}
	messageChan := make([]chan *kafka.Message, numChan)
	for i := 0; i < numChan; i++ {
		messageChan[i] = make(chan *kafka.Message, 1000)
	}
	go c.processMessage(messageChan, processFn)
	running := true
	for running == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			c.Logger.Infof("Caught signal %v: terminating\n", sig)
			running = false
		default:
			message, err := c.Consumer.ReadMessage(time.Second)
			if err != nil {
				newErr := err.(kafka.Error)
				if newErr.Code() != kafka.ErrTimedOut {
					c.Logger.Errorf("read message with error %v", err)
				}
				continue
			}
			messageChan[message.TopicPartition.Partition] <- message
		}

	}
	for {
		if _, err := c.Consumer.Commit(); err != nil {
			c.Logger.Errorf("commit failed due to %v, will retry soon", err)
			time.Sleep(1 * time.Second)
		}
		break
	}
	c.Logger.Info("Shutdown kafka consumer")
}
