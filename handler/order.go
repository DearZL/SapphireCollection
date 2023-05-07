package handler

import (
	"P/model"
	"P/resp"
	"P/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"time"
)

type OrderHandler struct {
	OrderSrvI service.OrderServiceInterface
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "订单创建失败",
		Data: nil,
	}
	order := &model.Order{
		OrderNum:    "20230001",
		SellerId:    "seller123",
		BuyerId:     "buyer456",
		OrderAmount: 100.0,
	}
	com1 := &model.Commodity{
		Model:        gorm.Model{ID: 3},
		Image:        "ds",
		Price:        222,
		Name:         "2131",
		OfferingDate: time.Now(),
		Status:       false,
		Number:       23,
	}
	com2 := &model.Commodity{
		Model:        gorm.Model{ID: 4},
		Image:        "ds",
		Price:        222,
		Name:         "2131",
		OfferingDate: time.Now(),
		Status:       false,
		Number:       23,
	}
	var commodities model.Commodities
	commodities.Commodities = append(commodities.Commodities, com1)
	commodities.Commodities = append(commodities.Commodities, com2)
	reOrder := order.ToRespOrder()
	reOrder.Commodities = commodities.ToRespCommodities()
	err := h.OrderSrvI.CreateOrder(order, commodities.Commodities)
	if err != nil {
		log.Println(err.Error())
		return
	}
	entity.Data = reOrder
	c.JSON(200, gin.H{"entity": entity})
}
