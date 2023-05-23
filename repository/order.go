package repository

import (
	"P/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type OrderRepoInterface interface {
	GetDB() *gorm.DB
	AddOrder(order *model.Order, tx ...*gorm.DB) error
	FindOrderByOrderNum(order *model.Order) error
	FindOrderWithComByOrderNum(order *model.Order) error
	UpdateOrderStatus(order *model.Order, status int, tx ...*gorm.DB) error
	DelOrder(order *model.Order, tx ...*gorm.DB) error
}
type OrderRepository struct {
	DB *gorm.DB
}

func (repo *OrderRepository) GetDB() *gorm.DB {
	db := repo.DB
	return db
}

func (repo *OrderRepository) AddOrder(order *model.Order, tx ...*gorm.DB) error {
	db := repo.DB
	if len(tx) != 0 {
		db = tx[0]
	}
	result := db.Create(order)
	fmt.Println(result.RowsAffected)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *OrderRepository) FindOrderByOrderNum(order *model.Order) error {
	result := repo.DB.Where("order_num=?", order.OrderNum).Find(order)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("未找到记录")
	}
	return nil
}

func (repo *OrderRepository) FindOrderWithComByOrderNum(order *model.Order) error {
	result := repo.DB.Preload("Commodities").Where("order_num=?", order.OrderNum).Find(order)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("未找到记录")
	}
	return nil
}

func (repo *OrderRepository) UpdateOrderStatus(order *model.Order, status int, tx ...*gorm.DB) error {
	db := repo.DB
	if len(tx) != 0 {
		db = tx[0]
	}
	err := db.Model(&model.Order{}).Where("order_num=?", order.OrderNum).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}
func (repo *OrderRepository) DelOrder(order *model.Order, tx ...*gorm.DB) error {
	db := repo.DB
	if len(tx) != 0 {
		db = tx[0]
	}
	err := db.Delete(order).Error
	if err != nil {
		return err
	}
	return nil
}
