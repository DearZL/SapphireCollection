package handler

import (
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
	com1 := &model.Commodity{
		Image:        "0f4bd884-dc9c-4cf9-b59e-7d5958fec3dd.jpg",
		Price:        222,
		Hash:         "231412412",
		Name:         "2131",
		OfferingDate: time.Now(),
		Status:       false,
		Number:       23,
	}
	com2 := &model.Commodity{
		Image:        "QQ图片20230210175331.jpg",
		Price:        241241222,
		Name:         "qweqwr",
		Hash:         "23124124",
		OfferingDate: time.Now(),
		Status:       false,
		Number:       23,
	}
	cs := &model.Commodities{}
	cs.Commodities = append(cs.Commodities, com1)
	cs.Commodities = append(cs.Commodities, com2)
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
