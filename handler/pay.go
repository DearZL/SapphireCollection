package handler

import (
	"P/model"
	"P/service"
	"github.com/gin-gonic/gin"
	"log"
)

type PayHandler struct {
	PaySrvI   service.PayServiceInterface
	OrderSrvI service.OrderServiceInterface
}

func (h *PayHandler) PayStatus(c *gin.Context) {
	order := &model.Order{
		OrderNum: c.PostForm("orderNum"),
	}
	err := h.OrderSrvI.FindOrderWithCom(order)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = h.PaySrvI.FindPayStatus(order)
	if err != nil {
		return
	}
	c.String(200, "OK")
	return
}
