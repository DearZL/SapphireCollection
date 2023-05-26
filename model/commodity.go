package model

import (
	"P/resp"
	"gorm.io/gorm"
	"time"
)

type Commodity struct {
	gorm.Model   `json:"gorm.Model"`
	Hash         string    `json:"hash" gorm:"unique"`
	Image        string    `json:"image"`
	Price        float64   `json:"price"`
	Name         string    `json:"name" gorm:"index"`
	Status       bool      `json:"status"`          //状态
	Number       int       `json:"number"`          //一组商品内的序号
	Amount       int       `json:"amount" gorm:"-"` //数量
	OfferingDate time.Time `json:"offeringDate"`    //发售日期
	OrderNum     string    `json:"orderNum"`        //订单编号
	UserId       string    `json:"userId"`
}
type Commodities struct {
	Commodities []*Commodity
}

func (c *Commodity) ToRespCommodity() *resp.Commodity {
	re := &resp.Commodity{
		Hash:     c.Hash,
		Image:    c.Image,
		Name:     c.Name,
		Price:    c.Price,
		Amount:   c.Amount,
		OrderNum: c.OrderNum,
	}
	return re
}

func (c *Commodity) ToUserCommodity() *resp.UserCommodity {
	re := &resp.UserCommodity{
		Hash:  c.Hash,
		Image: c.Image,
		Name:  c.Name,
	}
	return re
}

func (cs *Commodities) ToRespCommodities() resp.Commodities {
	var re resp.Commodities
	for _, c := range cs.Commodities {
		re.Commodities = append(re.Commodities, c.ToRespCommodity())
	}
	return re
}

func (cs *Commodities) ToRespUserCommodities() resp.UserCommodities {
	var re resp.UserCommodities
	for _, c := range cs.Commodities {
		re.UserCommodities = append(re.UserCommodities, c.ToUserCommodity())
	}
	return re
}
