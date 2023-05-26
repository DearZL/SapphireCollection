package service

import (
	"P/model"
	"P/repository"
	"errors"
	"log"
)

type CommodityService struct {
	CommodityRepo repository.CommodityRepoInterface
	BlockSrvI     BlockServiceInterface
}
type CommodityServiceInterface interface {
	AddCommodities(cs *model.Commodities) error
	FindCommoditiesByName(name string) (*model.Commodities, error)
	EditComSStatusByName(name string, status bool) error
}

func (srv *CommodityService) AddCommodities(cs *model.Commodities) error {
	tx := srv.BlockSrvI.GetDB().Begin()
	err := srv.BlockSrvI.AddBlock(cs, "0", tx)
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

func (srv *CommodityService) FindCommoditiesByName(name string) (*model.Commodities, error) {
	comS, err := srv.CommodityRepo.FindComSByName(name)
	if err != nil {
		return nil, err
	}
	cs := &model.Commodities{}
	cs.Commodities = comS
	return cs, nil
}

func (srv *CommodityService) EditComSStatusByName(name string, status bool) error {
	_, err := srv.FindCommoditiesByName(name)
	if err != nil {
		return err
	}
	err = srv.CommodityRepo.UpdateCommoditiesStatusByName(status, name)
	if err != nil {
		log.Println(err.Error())
		return errors.New("更新商品状态失败")
	}
	return nil
}
