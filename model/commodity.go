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
	Price        float32   `json:"price"`
	Name         string    `json:"name" gorm:"index"`
	Status       bool      `json:"status"`
	Number       int       `json:"number"`
	Amount       int       `json:"amount" gorm:"-"`
	OfferingDate time.Time `json:"offeringDate"`
	OrderNum     string    `json:"orderNum"`
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
