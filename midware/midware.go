package midware

import (
	"P/global"
	"P/resp"
	"P/utils"
	"crypto/sha256"
	"fmt"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"time"
)

var Store redis.Store

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization,x-requested-with,ClientId,FileType,token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,Authorization,ClientId")
		c.Header("Access-Control-Allow-Credentials", "true")
		method := c.Request.Method
		if method == "OPTIONS" {
			c.AbortWithStatus(200)
		}
		c.Next()
		return
	}
}

func GClientId() gin.HandlerFunc {
	return func(c *gin.Context) {
		entity := resp.EntityA{
			Code: 500,
			Msg:  "Error",
			Data: nil,
		}
		id := c.Request.Header.Get("ClientId")
		if id == "" {
			idx := sha256.Sum256([]byte(time.Now().String()))
			id = fmt.Sprintf("%x", idx[:])
			c.Header("ClientId", id)
			entity.SetCodeAndMsg(200, "")
			c.JSON(200, gin.H{"entity": entity})
			c.Next()
			return
		}
		c.Header("ClientId", id)
		entity.SetCodeAndMsg(200, "")
		c.JSON(200, gin.H{"entity": entity})
		c.Next()
		return
	}
}

func RequestFreq(freq int, t int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		entity := resp.EntityA{
			Code: 500,
			Msg:  "Error,No ClientId!",
			Data: nil,
		}
		id := c.Request.Header.Get("ClientId")
		if id == "" || len(id) != 64 {
			c.JSON(500, gin.H{"entity": entity})
			c.Abort()
			return
		}
		key := id + c.Request.RequestURI
		num, _ := strconv.Atoi(global.Redis1.Get(key).Val())
		if num >= freq {
			entity.SetCodeAndMsg(500, "您的请求过快,休息一下吧~")
			global.Redis1.Set(key, freq, time.Duration(t)*time.Second)
			c.JSON(200, gin.H{"entity": entity})
			c.Abort()
			return
		}
		global.Redis1.Set(key, num+1, time.Duration(t)*time.Second)
		c.Next()
		return
	}
}

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		entity := resp.EntityA{
			Code: 500,
			Msg:  "您的登录状态已失效,请重新登录！",
			Data: nil,
		}
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			log.Println("未携带Token")
			entity.SetCodeAndMsg(500, "非法请求")
			c.JSON(200, gin.H{"entity": entity})
			c.Abort()
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
		var d time.Duration
		if u.UserGroup == "Admin" {
			d = time.Duration(viper.GetInt64("loginTimeout.admin")) * time.Minute
			u.ExpiresAt = jwt.NewNumericDate(time.Now().Add(d))
			token, err = u.GToken()
			if err != nil {
				log.Println(err.Error())
				c.Abort()
				return
			}
		} else {
			if time.Now().Add(1*time.Hour).Compare(u.ExpiresAt.Time) == 1 {
				d = time.Duration(viper.GetInt64("loginTimeout.user")) * time.Hour
				u.ExpiresAt = jwt.NewNumericDate(time.Now().Add(d))
				token, err = u.GToken()
				if err != nil {
					log.Println(err.Error())
					c.Abort()
					return
				}
			}
		}
		global.Redis.Set(u.UserId, token, d)
		c.Header("Authorization", token)
		c.Set("token", token)
		c.Set("userId", u.UserId)
		c.Set("userGroup", u.UserGroup)
		c.Next()
		return
	}
}

func init() {
	var err error
	url := viper.GetString("redis.host") + ":" + viper.GetString("redis.port")
	password := viper.GetString("redis.password")
	Store, err = redis.NewStoreWithDB(2000, "tcp", url, password, "1", []byte(viper.GetString("storeKey")))
	if err != nil {
		panic(err.Error())
		return
	}
}
