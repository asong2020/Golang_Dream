package main

import (
	"asong.cloud/Golang_Dream/code_demo/kafka_demo/common"
	"asong.cloud/Golang_Dream/code_demo/kafka_demo/consumer"
	"asong.cloud/Golang_Dream/code_demo/kafka_demo/producer"
)

func main()  {
	topic := "asong_kafka_test"
	consm := consumer.NewConsumer(common.Group, []string{topic},&consumer.EventHandler{})
	defer consm.Stop()
	go consm.Consume() // 异步消费

	pro := producer.NewEventProducer(topic)

	for i := 0; i < 100; i++{
		msg := &common.KafkaMsg{
			ID: uint64(i),
			Detail: "asong就是玩",
		}
		pro.Producer(msg)
	}
	select {
	}
}
