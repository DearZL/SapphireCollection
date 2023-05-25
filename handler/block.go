package handler

import (
	"P/model"
	"P/resp"
	"P/service"
	"github.com/gin-gonic/gin"
)

type BlockHandler struct {
	BlockSrvI service.BlockServiceInterface
}

func (h *BlockHandler) BlocksInfo(c *gin.Context) {
	entity := resp.EntityA{
		Code:  500,
		Msg:   "查询失败",
		Token: "",
		Data:  nil,
	}
	if c.Param("chainId") == "" {
		entity.SetCodeAndMsg(500, "参数错误")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	chain := &model.Blockchain{Id: c.Param("chainId")}
	err := h.BlockSrvI.GetBlockChain(chain)
	if err != nil {
		entity.SetCodeAndMsg(500, err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.Data = chain.Blocks
	entity.SetCodeAndMsg(200, "查询成功")
	entity.SetEntityAndHeaderToken(c)
	c.JSON(200, gin.H{"entity": entity})
	return
}
