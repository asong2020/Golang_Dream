package common

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
)

const (
	VERSION = "2.6.0"
	ADDR = "localhost:9092"
	Group = "asong"
)

func NewAsyncProducer() sarama.AsyncProducer {
	cfg := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(VERSION)
	if err != nil{
		log.Fatal("NewAsyncProducer Parse kafka version failed", err.Error())
		return nil
	}
	cfg.Version = version
	cfg.Producer.RequiredAcks = sarama.WaitForAll // 三种模式任君选择
	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Retry.Max = 3
	cfg.Producer.Retry.Backoff = 100 * time.Millisecond
	cli, err := sarama.NewAsyncProducer([]string{ADDR}, cfg)
	if err != nil{
		log.Fatal("NewAsyncProducer failed", err.Error())
		return nil
	}
	return cli
}

// 定义消息结构
type KafkaMsg struct {
	Detail string
	ID uint64
}

func NewConsumerGroup(group string) sarama.ConsumerGroup {
	cfg := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(VERSION)
	if err != nil{
		log.Fatal("NewConsumerGroup Parse kafka version failed", err.Error())
		return nil
	}

	cfg.Version = version
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Offsets.Retry.Max = 3
	cfg.Consumer.Offsets.AutoCommit.Enable = true // 开启自动提交，需要手动调用MarkMessage才有效
	cfg.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 间隔
	client, err := sarama.NewConsumerGroup([]string{ADDR}, group, cfg)
	if err != nil {
		log.Fatal("NewConsumerGroup failed", err.Error())
	}
	return client
}