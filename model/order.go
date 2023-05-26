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
	CommodityNumber int          `json:"commodityNumber"`       //商品个数
	OrderAmount     float64      `json:"orderAmount"`           //订单总金额
	OrderType       string       `json:"orderType"`             //订单类型:wallet||commodity
	Status          int          `json:"status" gorm:"NotNull"` //订单状态
}

// ToRespOrder 转换为响应订单类型
func (o *Order) ToRespOrder() *resp.Order {
	var comS Commodities
	comS.Commodities = o.Commodities
	reOrder := &resp.Order{
		OrderNum:     o.OrderNum,
		SellerId:     o.SellerId,
		BuyerId:      o.BuyerId,
		Commodities:  comS.ToRespCommodities(),
		CommodityNum: o.CommodityNumber,
		OrderAmount:  o.OrderAmount,
		OrderTime:    o.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return reOrder
}
