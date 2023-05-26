package router

import (
	"P/global"
	"P/midware"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)
	gin.SetMode(viper.GetString("app.mode"))
	r.Use(sessions.Sessions("LD_S", midware.Store)).Use(midware.Cors())

	pay := r.Group("/api/pay")
	{
		pay.POST("/payStatus", global.PayHandler.PayStatus)
		pay.GET("/success", global.PayHandler.JumpQueryPayStatus)
	}

	common := r.Group("/api/common")
	{
		common.GET("/GClientId", midware.GClientId())
		common.POST("/upload", midware.CheckToken(), global.CommonHandler.Upload)
		common.GET("/download/:filename", global.CommonHandler.Download)
		common.GET("/testToken", midware.CheckToken(), func(c *gin.Context) {
			userId, _ := c.Get("userId")
			fmt.Println(userId.(string) == "asdwqeq")
			c.String(200, "Token合法")
		})
	}

	user := r.Group("/api/user").Use(midware.RequestFreq(5, 3))
	{
		user.GET("/checkEmail/:email", global.UserHandler.CheckEmail)
		user.POST("/userReg/:code", global.UserHandler.UserReg)
		user.POST("/sendCode", global.UserHandler.SendCode)
		user.POST("/login", global.UserHandler.UserLogin)
		user.POST("/editUserPw", global.UserHandler.EditUserPw)
		user.GET("/logout", midware.CheckToken(), global.UserHandler.Logout)
		user.GET("/userInfo", midware.CheckToken(), global.UserHandler.UserInfo)
		user.POST("/editUser", midware.CheckToken(), global.UserHandler.EditUserInfo)
		user.POST("/emailVerify", midware.CheckToken(), global.UserHandler.EmailVerify)
		user.POST("/editUserEmail", midware.CheckToken(), global.UserHandler.EditUserEmail)
		user.POST("/editUserAdmin", midware.CheckToken(), global.UserHandler.AdminEditUser)
		user.POST("/editUserStatus", midware.CheckToken(), global.UserHandler.EditUserStatus)
		user.POST("/wallet", midware.CheckToken())
	}
	block := r.Group("/api/block")
	{
		block.GET("/findChain/:chainId", global.BlockHandler.BlocksInfo)
	}

	commodity := r.Group("/api/commodity")
	{
		commodity.POST("/addCommodity", global.CommodityHandler.AddCommodities)
		commodity.POST("/editComSStatus", global.CommodityHandler.EditComSStatus)
	}

	order := r.Group("/api/order").Use(midware.CheckToken())
	{
		order.POST("/createOrder", global.OrderHandler.CreateOrder)
		order.POST("/createWalletOrder", global.OrderHandler.CreateWalletOrder)
		order.POST("/payOrder", global.OrderHandler.PayOrder)
		order.POST("/dropOrder", global.OrderHandler.DropOrder)
	}

	return r
}
