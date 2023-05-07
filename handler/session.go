package handler

import (
	"P/resp"
	"P/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	SessionSrvI service.SessionServiceInterface
}

func (h *SessionHandler) DropSession(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "清除session失败！",
		Data: nil,
	}
	session := sessions.Default(c)
	sid := session.ID()
	st, err := h.SessionSrvI.DropSession(sid)
	if err != nil {
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	if st == "OK" {
		entity.SetMsg("清除session成功！")
		entity.SetCode(200)
	}
	c.JSON(200, gin.H{"entity": entity})
	return
}
