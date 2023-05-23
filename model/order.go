package model

import (
	"P/resp"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model      `json:"gorm.Model" `
	OrderNum        string       `json:"orderNum" gorm:"unique;type:varchar(255)"`
	SellerId        string       `json:"sellerId"`
	BuyerId         string       `json:"buyerId"`
	Commodities     []*Commodity `json:"commodities" gorm:"foreignKey:OrderNum;references:OrderNum"`
	CommodityAmount int          `json:"commodityAmount"`
	OrderAmount     float32      `json:"orderAmount"`
	Status          int          `json:"status" gorm:"NotNull"`
}

func (o *Order) ToRespOrder() *resp.Order {
	var comS Commodities
	comS.Commodities = o.Commodities
	reOrder := &resp.Order{
		OrderNum:    o.OrderNum,
		SellerId:    o.SellerId,
		BuyerId:     o.BuyerId,
		Commodities: comS.ToRespCommodities(),
		OrderAmount: o.OrderAmount,
		OrderTime:   o.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return reOrder
}
