package consumer

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Shopify/sarama"

	"asong.cloud/Golang_Dream/code_demo/kafka_demo/common"
)

type EventHandler struct {

}

func (e EventHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (e EventHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (e EventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var data common.KafkaMsg
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			return errors.New("failed to unmarshal message err is " + err.Error())
		}
		// 操作数据，改用打印
		log.Print("consumerClaim data is ")

		// 处理消息成功后标记为处理, 然后会自动提交
		session.MarkMessage(msg,"")
	}
	return nil
}