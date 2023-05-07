package repository

import (
	"P/model"
	"gorm.io/gorm"
	"log"
)

type OrderRepoInterface interface {
	AddOrder(tx *gorm.DB, order *model.Order) error
	GetDB() *gorm.DB
}
type OrderRepository struct {
	DB *gorm.DB
}

func (re *OrderRepository) GetDB() *gorm.DB {
	db := re.DB
	return db
}
func (re *OrderRepository) AddOrder(tx *gorm.DB, order *model.Order) error {
	err := tx.Create(order).Error
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
