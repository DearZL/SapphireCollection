package service

import (
	"P/model"
	"P/repository"
)

type CommodityService struct {
	CommodityRepo repository.CommodityRepoInterface
	BlockSrv      BlockService
}
type CommodityServiceInterface interface {
	AddCommodities(cs []*model.Commodity) error
}

func (srv *CommodityService) AddCommodities(cs []*model.Commodity) error {
	return srv.CommodityRepo.AddCommodities(cs)
}
