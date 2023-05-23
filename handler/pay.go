package handler

import (
	"P/model"
	"P/resp"
	"P/service"
	"github.com/gin-gonic/gin"
	"log"
)

type PayHandler struct {
	PaySrvI   service.PayServiceInterface
	OrderSrvI service.OrderServiceInterface
}

func (h *PayHandler) PayStatus(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "没有此订单记录,请检查参数后重试",
		Data: nil,
	}
	if c.PostForm("orderNum") == "" {
		entity.SetCodeAndMsg(500, "参数错误！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	order := &model.Order{
		OrderNum: c.PostForm("orderNum"),
	}
	err := h.OrderSrvI.FindOrderWithCom(order)
	if err != nil {
		log.Println(err.Error())
		entity.SetCodeAndMsg(500, "查询订单出错")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	result, err := h.PaySrvI.FindPayStatus(order)
	if err != nil {
		entity.SetCodeAndMsg(500, "查询订单出错")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetEntityAndHeaderToken(c)
	entity.SetCodeAndMsg(200, "查询订单成功")
	entity.Data = result.Content
	return
}
