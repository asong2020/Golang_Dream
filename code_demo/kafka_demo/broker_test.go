package main

import (
	"fmt"
	"testing"

	"asong.cloud/Golang_Dream/code_demo/kafka_demo/common"
)

func TestNewAsyncProducer(t *testing.T) {
	cli := common.NewAsyncProducer()
	fmt.Println(cli)
}