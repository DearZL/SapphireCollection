package model

import (
	"P/resp"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model  `json:"gorm.Model" `
	OrderNum    string  `json:"orderNum"`
	SellerId    string  `json:"sellerId"`
	BuyerId     string  `json:"buyerId"`
	OrderAmount float32 `json:"orderAmount"`
}

func (o *Order) ToRespOrder() *resp.Order {
	reOrder := &resp.Order{
		OrderNum:    o.OrderNum,
		SellerId:    o.SellerId,
		BuyerId:     o.BuyerId,
		OrderAmount: o.OrderAmount,
	}
	return reOrder
}
