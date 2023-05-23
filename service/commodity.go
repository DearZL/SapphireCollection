package service

import (
	"P/model"
	"P/repository"
	"log"
)

type CommodityService struct {
	CommodityRepo repository.CommodityRepoInterface
	BlockSrv      *BlockService
}
type CommodityServiceInterface interface {
	AddCommodities(cs *model.Commodities) error
}

func (srv *CommodityService) AddCommodities(cs *model.Commodities) error {
	tx := srv.BlockSrv.BlockRepo.GetDB().Begin()

	err := srv.BlockSrv.AddBlock(cs, "0", tx)
	if err != nil {
		tx.Rollback()
		log.Println("rollback")
		return err
	}
	err = srv.CommodityRepo.AddCommodities(cs.Commodities, tx)
	if err != nil {
		tx.Rollback()
		log.Println("rollback")
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		log.Println("rollback")
		return err
	}
	return nil
}
