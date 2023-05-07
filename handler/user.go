package handler

import (
	"P/model"
	"P/resp"
	"P/service"
	"P/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"strconv"
	"time"
)

type UserHandler struct {
	UserSrvI service.UserSrvInterface
	Redis    *redis.Client
}

// CheckEmail 检查邮箱
func (h *UserHandler) CheckEmail(c *gin.Context) {
	//初始化响应结构体
	entity := resp.EntityA{
		Code: 500,
		Msg:  "查询失败",
		Data: nil,
	}
	u := model.User{Email: c.Param("email")}
	if u.Email == "" {
		entity.SetMsg("邮箱不能为空！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//查询数据库中是否已经有该邮箱
	r := h.UserSrvI.EmailExist(u)
	if r != nil {
		entity.SetMsg("该邮箱已被注册！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetMsg("该邮箱可用")
	entity.SetCode(200)
	c.JSON(200, gin.H{"entity": entity})
}

// SendCode 发送验证码
func (h *UserHandler) SendCode(c *gin.Context) {
	//定义响应结构体
	entity := resp.EntityA{
		Code: 500,
		Msg:  "验证码发送失败",
		Data: nil,
	}
	//初始化一个空User对象用于接收参数
	u := model.User{}
	//接收前端传来的email参数并放入user中的email字段
	u.Email = c.Param("email")
	if u.Email == "" {
		entity.SetMsg("邮箱不能为空")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//为本次会话上下文初始化一个session
	session := sessions.Default(c)
	//设置session的有效期为五分钟，及其他属性
	session.Options(sessions.Options{MaxAge: 5 * 60, HttpOnly: true, Path: "/"})
	//调用工具类生成随机数并转为string类型
	code := strconv.Itoa(utils.RandInt(123456, 999999))
	//将该验证码放入session中
	session.Set("code", code)
	//保存session
	err := session.Save()
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//调用service层发送验证码
	err = h.UserSrvI.SendCode(code, u)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetMsg("验证码发送成功")
	entity.SetCode(200)
	c.JSON(200, gin.H{"entity": entity})
}

func (h *UserHandler) AddUser(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "用户注册失败",
		Data: nil,
	}
	//初始化一个user对象接收前端参数
	u := model.User{}
	//将前端传来的json序列化成user对象
	err := c.ShouldBindJSON(&u)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//如果不含有email值则非法
	if u.Email == "" || u.Password == "" {
		entity.SetMsg("邮箱或密码不能为空!")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//创建session
	session := sessions.Default(c)
	//取出该上下文中session中code的值
	code := session.Get("code")
	//若前端没有传回code则返回
	if c.Param("code") == "" || code == nil {
		entity.Msg = "验证码不能为空！"
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//将参数中的code和session中取出的code作比较，其中code.(string)是将code（类型interface{}）断言为string类型
	if c.Param("code") != code.(string) {
		entity.Msg = "验证码错误！"
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//验证成功后删除session中的code
	session.Delete("code")
	//保存session
	sessionErr := session.Save()
	if sessionErr != nil {
		log.Println(sessionErr.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//调用service层添加该用户
	user, err := h.UserSrvI.Add(u)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetCode(200)
	entity.SetMsg("注册成功！")
	//将user转换为响应user类型避免暴露过多字段
	entity.Data = user.ToRespUser()
	c.JSON(200, gin.H{"entity": entity})
}

func (h *UserHandler) UserLogin(c *gin.Context) {
	var (
		t time.Time
		d time.Duration
	)
	entity := resp.EntityA{
		Code: 500,
		Msg:  "登陆失败",
		Data: nil,
	}
	u := &model.User{}
	err := c.ShouldBindJSON(u)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	if u.Email == "" || u.Password == "" {
		entity.SetMsg("邮箱密码不能为空！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//向数据库查询该邮箱是否存在
	r := h.UserSrvI.EmailExist(*u)
	if r == nil {
		entity.SetMsg("该用户不存在，请检查邮箱或注册！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	if r.Status != false {
		entity.SetMsg("该用户已被禁用")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//查询数据库中的用户密码
	user, err := h.UserSrvI.UserLogin(*u)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//比对密码
	if u.Password != user.Password {
		entity.SetMsg("密码或邮箱错误，请检查后重新输入！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//鉴别身份以生成不同有效期的token
	if user.UserId == "Admin" {
		d = 10 * time.Second
	} else {
		d = 72 * time.Hour
	}
	t = time.Now().Add(d)
	//初始化一个claims结构体用以生成token
	claims := utils.UserClaims{
		UserId: user.UserId,
		//生成随机数防止撞车
		Num: strconv.Itoa(utils.RandInt(123456, 999999)),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t),
		},
	}
	//调用工具类生成token
	token, err := claims.GToken()
	if err != nil {
		log.Println(err.Error())
		return
	}
	session := sessions.Default(c)
	session.Options(sessions.Options{MaxAge: 3600 * 72, HttpOnly: true, Path: "/"})
	//将用户信息存入session
	session.Set("user", user)
	//保存
	err = session.Save()
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetMsg("登陆成功")
	entity.SetCode(200)
	//设置响应数据
	entity.Token = token
	entity.Data = user.ToRespUser()
	//向redis中写入userID和token
	h.Redis.Set(user.UserId, token, d)
	//设置自定义响应头
	c.Header("Authorization", token)
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *UserHandler) UserInfo(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "Error",
		Data: nil,
	}
	session := sessions.Default(c)
	userInfo := session.Get("user")
	if userInfo == nil {
		entity.SetMsg("您的会话已过期，请重新登录!")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	u := model.User{
		UserId: userInfo.(*model.User).UserId,
	}
	user, err := h.UserSrvI.UserInfo(u)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetCode(200)
	entity.SetMsg("OK")
	entity.Data = user.ToRespUser()
	c.JSON(200, gin.H{"entity": entity})
}

func (h *UserHandler) Logout(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "Error",
		Data: nil,
	}
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetMsg("退出成功！")
	entity.SetCode(200)
	c.Header("Authorization", "")
	c.JSON(200, gin.H{"entity": entity})
}

func (h *UserHandler) EditUser(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "参数错误！",
		Data: nil,
	}
	user := &model.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	u := &utils.UserClaims{}
	token := c.Request.Header.Get("Authorization")
	s, err := u.PToken(token)
	u.MapClaims(s)
	if u.UserId != "Admin" && user.UserId != u.UserId {
		entity := resp.EntityA{
			Code: 500,
			Msg:  "非法请求！",
			Data: nil,
		}
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	err = h.UserSrvI.EditUser(*user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetMsg("更改信息成功!")
	entity.SetCode(200)
	Token, _ := c.Get("token")
	entity.Token = Token.(string)
	c.JSON(200, gin.H{"entity": entity})

}
