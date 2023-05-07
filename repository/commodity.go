package repository

import (
	"P/model"
	"gorm.io/gorm"
	"log"
)

type CommodityRepository struct {
	DB *gorm.DB
}
type CommodityRepoInterface interface {
	AddCommodityWithOrder(tx *gorm.DB, Commodities []*model.Commodity) error
}

func (repo *CommodityRepository) AddCommodityWithOrder(tx *gorm.DB, Commodities []*model.Commodity) error {
	err := tx.Create(Commodities).Error
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
