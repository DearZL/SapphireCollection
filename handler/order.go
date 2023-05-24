package handler

import (
	"P/model"
	"P/resp"
	"P/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
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
	com1 := &model.Commodity{
		Name: "2132",
	}
	order := &model.Order{
		SellerId: "seller123",
		BuyerId:  "buyer456",
		//商品总数量
		CommodityAmount: 2,
		Commodities:     []*model.Commodity{com1},
	}
	order, err := h.OrderSrvI.CreateOrder(order, com1)
	if err != nil {
		log.Println(err.Error())
		entity.SetCodeAndMsg(500, err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	respOrder := order.ToRespOrder()
	entity.SetEntityAndHeaderToken(c)
	entity.SetCodeAndMsg(200, "订单创建成功!请于"+viper.GetString("order.timeout")+"分钟内支付！")
	entity.Data = respOrder
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *OrderHandler) PayOrder(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "生成支付链接失败!",
		Data: nil,
	}
	if c.PostForm("orderNum") == "" {
		entity.SetCodeAndMsg(500, "参数错误!")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	order := &model.Order{OrderNum: c.PostForm("orderNum")}
	err := h.OrderSrvI.FindOrder(order)
	if err != nil {
		log.Println(err.Error())
		entity.SetCodeAndMsg(500, "未找到订单信息!")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	payUrl, err := h.OrderSrvI.PayOrder(order)
	if err != nil {
		entity.SetCodeAndMsg(500, err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.Data = payUrl.String()
	entity.SetEntityAndHeaderToken(c)
	entity.SetCodeAndMsg(200, "生成支付链接成功!")
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *OrderHandler) DropOrder(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "订单取消失败",
		Data: nil,
	}
	if c.PostForm("orderNum") == "" {
		entity.SetCodeAndMsg(500, "参数错误")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	order := &model.Order{
		OrderNum: c.PostForm("orderNum"),
	}

	err := h.OrderSrvI.DropOrder(order)
	if err != nil {
		log.Println(err.Error())
		entity.SetCodeAndMsg(500, err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetEntityAndHeaderToken(c)
	entity.SetCodeAndMsg(200, "订单取消成功")
	c.JSON(200, gin.H{"entity": entity})
	return
}
