package common

import (
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic/v7"

	"asong.cloud/Golang_Dream/code_demo/go-elastic-asong/config"
)

func NewEsClient(conf *config.ServerConfig) *elastic.Client {
	url := fmt.Sprintf("http://%s:%d", conf.Elastic.Host, conf.Elastic.Port)
	client, err := elastic.NewClient(
		//elastic 服务地址
		elastic.SetURL(url),
		// 设置错误日志输出
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {
		log.Fatalln("Failed to create elastic client")
	}
	return client
}
