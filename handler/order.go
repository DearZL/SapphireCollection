package handler

import (
	"P/enum"
	"P/model"
	"P/resp"
	"P/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"strconv"
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
	order := &model.Order{}
	userId, exist := c.Get("userId")
	if err := c.ShouldBindJSON(&order.Commodities); err != nil || !exist {
		log.Println(err)
		entity.SetCodeAndMsg(500, "参数错误")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	order.OrderType = enum.OrderTypeCommodity
	order.BuyerId = userId.(string)
	//商品总数量
	for _, c := range order.Commodities {
		order.CommodityNumber = order.CommodityNumber + c.Amount
	}
	if order.CommodityNumber <= 0 {
		entity.SetCodeAndMsg(500, "非法参数")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	err := h.OrderSrvI.CreateOrder(order)
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

func (h *OrderHandler) CreateWalletOrder(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "参数不存在,请求错误",
		Data: nil,
	}
	userId, exist := c.Get("userId")
	if !exist || c.PostForm("money") == "" {
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	order := &model.Order{}
	//断言userId类型及补充订单
	order.BuyerId = userId.(string)
	order.SellerId = "Admin"
	order.OrderType = enum.OrderTypeWallet
	amount, err := strconv.ParseFloat(c.PostForm("money"), 64)
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	if amount <= viper.GetFloat64("wallet.minTopUpAmount") {
		entity.SetCodeAndMsg(500, "充值金额不能小于"+viper.GetString("wallet.minTopUpAmount")+"元")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	order.OrderAmount = amount
	err = h.OrderSrvI.CreateOrder(order)
	if err != nil {
		log.Println(err)
		entity.SetCodeAndMsg(500, "创建订单失败")
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
