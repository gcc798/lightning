// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/gcc798/lightning/application/sys-api/internal/config"
	"github.com/gcc798/lightning/application/sys-rpc/client/sysservice"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	SysRpcClient sysservice.SysService
	Redis        *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.Db,
	})

	return &ServiceContext{
		Config:       c,
		SysRpcClient: sysservice.NewSysService(zrpc.MustNewClient(c.SysRpc)),
		Redis:        rdb,
	}
}
