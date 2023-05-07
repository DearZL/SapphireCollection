package router

import (
	"P/global"
	"P/midware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)
	gin.SetMode("debug")
	r.Use(sessions.Sessions("LD_S", midware.Store)).Use(midware.Cors())

	pay := r.Group("/api/pay")
	{
		pay.GET("/pay", global.PayHandler.Pay)
	}

	common := r.Group("/api/common").Use(midware.CheckToken())
	{
		common.POST("/upload", global.CommonHandler.Upload)
		common.GET("/download/:filename", global.CommonHandler.Download)
		common.GET("/testToken", func(c *gin.Context) {
			c.String(200, "Token合法")
		})
	}

	session := r.Group("/api/session").Use(midware.CheckToken())
	{
		session.GET("/DropSession", global.SessionHandler.DropSession)
	}

	user := r.Group("/api/user")
	{
		user.GET("/checkEmail/:email", global.UserHandler.CheckEmail)
		user.POST("/add/:code", global.UserHandler.AddUser)
		user.GET("/sendCode/:email", global.UserHandler.SendCode)
		user.POST("/login", global.UserHandler.UserLogin)
		user.GET("/logout", global.UserHandler.Logout)
		user.GET("/userInfo", midware.CheckToken(), global.UserHandler.UserInfo)
		user.POST("/editUser", midware.CheckToken(), global.UserHandler.EditUser)
	}

	order := r.Group("/api/order")
	{
		order.GET("/create", global.OrderHandler.CreateOrder)
	}

	return r
}
