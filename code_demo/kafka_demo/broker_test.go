package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"asong.cloud/Golang_Dream/code_demo/kafka_demo/common"
)

func TestNewAsyncProducer(t *testing.T) {
	cli := common.NewAsyncProducer()
	if assert.Equal(t,nil,cli) {
		fmt.Println("create cli failed")
	}
}