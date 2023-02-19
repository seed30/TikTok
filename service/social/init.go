package social

import (
	"github.com/seed30/TikTok/dal/dao"
	"github.com/seed30/TikTok/rpc"
)

type Service struct {
	dao *dao.Dao
	rpc *rpc.RPC
}

func NewService(rpc *rpc.RPC) *Service {
	return &Service{
		dao: dao.NewDao(),
		rpc: rpc,
	}
}
