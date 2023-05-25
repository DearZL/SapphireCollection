package handler

import (
	"P/enum"
	"P/model"
	"P/resp"
	"P/service"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type CommodityHandler struct {
	CommoditySrvI service.CommodityServiceInterface
}

func (h *CommodityHandler) AddCommodities(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "参数错误!",
		Data: nil,
	}
	com1 := &model.Commodity{}
	if err := c.ShouldBindJSON(com1); err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	com1.OfferingDate = time.Now()
	com1.Status = enum.CommodityDisabled
	com1.UserId = "Admin"
	cs := &model.Commodities{}
	for i := 1; i <= com1.Number; i++ {
		com := &model.Commodity{
			Number:       i,
			Name:         com1.Name,
			Image:        com1.Image,
			Price:        com1.Price,
			UserId:       com1.UserId,
			Status:       com1.Status,
			OfferingDate: com1.OfferingDate,
		}
		cs.Commodities = append(cs.Commodities, com)
	}
	err := h.CommoditySrvI.AddCommodities(cs)
	if err != nil {
		for _, c := range cs.Commodities {
			log.Println("商品添加失败!Name:", c.Name, "Hash:", c.Hash)
		}
		log.Println(err.Error())
		entity.SetCodeAndMsg(500, "添加商品失败")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetEntityAndHeaderToken(c)
	entity.SetCodeAndMsg(200, "添加商品成功!")
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *CommodityHandler) EditComSStatus(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "参数错误!",
		Data: nil,
	}
	if c.PostForm("status") == "" || c.PostForm("name") == "" {
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	err := h.CommoditySrvI.EditComSStatusByName(c.PostForm("name"), c.PostForm("status") == "true")
	if err != nil {
		entity.SetCodeAndMsg(500, err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetCodeAndMsg(200, "更改成功")
	entity.SetEntityAndHeaderToken(c)
	c.JSON(200, gin.H{"entity": entity})
	return
}
