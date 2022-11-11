package mq

import (
	"context"
	"ddd-demo/common"
	"encoding/json"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

// KafkaClient Kafka 客户端数据结构
type KafkaClient struct {
	brokers []string
}

// GetConsumerClint 获取消费者客户端
func (k *KafkaClient) GetConsumerClint(group string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumer, err := sarama.NewConsumerGroup(k.brokers, group, config)
	return consumer, err
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
	} else {
		client.PauseAll()
	}

	*isPaused = !*isPaused
}

// Consume 消费
func (k *KafkaClient) Consume(topic, group string, callback func(interface{})) error {
	keepRunning := true

	consumer := Consumer{
		ready:    make(chan bool),
		callback: callback,
	}

	client, err := k.GetConsumerClint(group)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, strings.Split(topic, ","), &consumer); err != nil {
				return
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			keepRunning = false
		case <-sigterm:
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		return err
	}

	return nil
}

// GetProducerClient 获取生产者客户端
func (k *KafkaClient) GetProducerClient() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner
	producer, err := sarama.NewSyncProducer(k.brokers, config)
	return producer, err
}

// Produce 生产
func (k *KafkaClient) Produce(topic string, message map[string]interface{}) error {
	client, err := k.GetProducerClient()
	defer func() {
		_ = client.Close()
	}()
	if err != nil {
		return err
	}

	value, err := json.Marshal(message)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(common.GetMd5(value)),
		Value:     sarama.ByteEncoder(value),
		Timestamp: time.Now(),
	}
	_, _, err = client.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

var (
	defaultKafkaClientOnce sync.Once
	defaultKafkaClient     *KafkaClient
)

// NewKafkaClient 创建 KafkaClient 对象
func NewKafkaClient(brokers []string) *KafkaClient {
	defaultKafkaClientOnce.Do(func() {
		defaultKafkaClient = &KafkaClient{
			brokers: brokers,
		}
	})
	return defaultKafkaClient
}

// Consumer 消费者结构体
type Consumer struct {
	ready    chan bool
	callback func(interface{})
}

// Message 消息结构
type Message struct {
	Topic     string
	Value     string
	Timestamp time.Time
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		c.callback(&Message{
			Topic:     string(message.Topic),
			Value:     string(message.Value),
			Timestamp: message.Timestamp,
		})
		session.MarkMessage(message, "")
	}

	return nil
}
