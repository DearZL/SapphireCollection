package handler

import (
	"P/service"
	"github.com/gin-gonic/gin"
)

type BlockHandler struct {
	BlockSrvI service.BlockServiceInterface
}

func (h *BlockHandler) BlocksInfo(c *gin.Context) {

}
