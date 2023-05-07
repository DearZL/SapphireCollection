package service

import (
	"P/model"
	"P/repository"
	"log"
)

type OrderServiceInterface interface {
	CreateOrder(order *model.Order, commodities []*model.Commodity) error
}

type OrderService struct {
	OrderRepo     repository.OrderRepoInterface
	CommodityRepo repository.CommodityRepoInterface
}

func (o *OrderService) CreateOrder(order *model.Order, commodities []*model.Commodity) error {
	db := o.OrderRepo.GetDB()
	tx := db.Begin()
	err := o.OrderRepo.AddOrder(tx, order)
	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}
	for _, commodity := range commodities {
		commodity.OrderID = order.OrderNum
	}
	err = o.CommodityRepo.AddCommodityWithOrder(tx, commodities)
	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
