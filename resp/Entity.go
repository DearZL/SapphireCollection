package resp

import (
	"github.com/gin-gonic/gin"
)

type EntityA struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Token string      `json:"token"`
	Data  interface{} `json:"data"`
}
type EntityB struct {
	Total     int `json:"total"`
	TotalPage int `json:"totalPage"`
	EntityA
}

func (e *EntityA) SetCodeAndMsg(code int, s string) {
	e.Code = code
	e.Msg = s
}

func (e *EntityA) SetEntityAndHeaderToken(c *gin.Context) {
	token, exists := c.Get("token")
	if exists {
		e.Token = token.(string)
		c.Header("Authorization", token.(string))
	}
}

func (e *EntityB) SetTotal(Total int, TotalPage int) {
	e.Total = Total
	e.TotalPage = TotalPage
}
