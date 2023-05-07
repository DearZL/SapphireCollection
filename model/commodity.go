package model

import (
	"P/resp"
	"gorm.io/gorm"
	"time"
)

type Commodity struct {
	gorm.Model   `json:"gorm.Model"`
	Hash         []byte    `json:"hash,omitempty"`
	Image        string    `json:"image,omitempty"`
	Price        float32   `json:"price,omitempty"`
	Name         string    `json:"name,omitempty"`
	Status       bool      `json:"status,omitempty"`
	Number       int       `json:"number,omitempty"`
	OfferingDate time.Time `json:"offeringDate,omitempty"`
	OrderID      string    `json:"orderID,omitempty"`
}
type Commodities struct {
	Commodities []*Commodity
}

func (c *Commodity) ToRespCommodity() *resp.Commodity {
	re := &resp.Commodity{
		Hash:    c.Hash,
		Image:   c.Image,
		Name:    c.Name,
		Price:   c.Price,
		OrderID: c.OrderID,
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
