package midware

import (
	"P/global"
	"P/resp"
	"P/utils"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"strconv"
	"time"
)

var Store redis.Store

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		method := c.Request.Method
		if method == "OPTIONS" {
			c.AbortWithStatus(200)
		}
		c.Next()
	}
}

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			t time.Time
			d time.Duration
		)
		entity := resp.EntityA{
			Code: 500,
			Msg:  "您的登录状态已失效,请重新登录！",
			Data: nil,
		}
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			log.Println("未携带Token")
			c.AbortWithStatus(404)
			return
		}
		u := &utils.UserClaims{}
		s, err := u.PToken(token)
		if err != nil {
			log.Println(err.Error())
			c.JSON(200, gin.H{"entity": entity})
			c.Abort()
			return
		}
		u.MapClaims(s)
		result := global.Redis.Get(u.UserId).Val()
		if result == "" || result != token {
			c.JSON(200, gin.H{"entity": entity})
			c.Abort()
			return
		}
		if u.UserId == "Admin" {
			d = 10 * time.Minute
			t = time.Now().Add(d)
			u.Num = strconv.Itoa(utils.RandInt(123456, 999999))
			u.ExpiresAt = jwt.NewNumericDate(t)
			token, err = u.GToken()
			if err != nil {
				log.Println(err.Error())
				c.Abort()
				return
			}
			global.Redis.Set(u.UserId, token, d)
			c.Header("Authorization", token)
			c.Set("token", token)
			c.Next()
		} else {
			if time.Now().Add(1*time.Hour).Compare(u.ExpiresAt.Time) == 1 {
				d = 72 * time.Hour
				t = time.Now().Add(d)
				u.ExpiresAt = jwt.NewNumericDate(t)
				token, err = u.GToken()
				if err != nil {
					log.Println(err.Error())
					c.Abort()
					return
				}
			}
			c.Header("Authorization", token)
			c.Set("token", token)
			c.Next()
		}
	}
}

func init() {
	var err error
	Store, err = redis.NewStoreWithDB(2000, "tcp", "192.168.124.6:6379", "", "1", []byte("Zhoul0722"))
	if err != nil {
		panic(err.Error())
		return
	}
}

//Store = gorm.NewStore(global.DB, true, []byte("ZhouL0722"))
