//+build wireinject

package wire

import (
	"github.com/google/wire"

	"asong.cloud/Golang_Dream/wire_cron_example/config"
	"asong.cloud/Golang_Dream/wire_cron_example/cron"
	"asong.cloud/Golang_Dream/wire_cron_example/cron/task"
	"asong.cloud/Golang_Dream/wire_cron_example/dao"
	"asong.cloud/Golang_Dream/wire_cron_example/service"
)

func InitializeCron(mysql *config.Mysql)  *cron.Cron{
	wire.Build(
		dao.NewClientDB,
		dao.NewUserDB,
		service.NewUserService,
		task.NewScanner,
		cron.NewCron,
		)
	return &cron.Cron{}
}