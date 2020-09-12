//+build wireinject

package wire

import (
	"github.com/google/wire"

	"asong.cloud/Golang_Dream/code_demo/go-elastic-asong/common"
	"asong.cloud/Golang_Dream/code_demo/go-elastic-asong/config"
	"asong.cloud/Golang_Dream/code_demo/go-elastic-asong/dao"
	"asong.cloud/Golang_Dream/code_demo/go-elastic-asong/handler"
	"asong.cloud/Golang_Dream/code_demo/go-elastic-asong/service"
)

func InitializeHandler(conf *config.ServerConfig) *handler.UserHandler{
	wire.Build(
		common.NewEsClient,
		common.NewRouterClient,

		dao.NewUserES,
		service.NewUserService,

		handler.NewUserHandler,
		)
	return &handler.UserHandler{}
}