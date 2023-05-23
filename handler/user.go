package handler

import (
	"P/model"
	"P/resp"
	"P/service"
	"P/utils"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type UserHandler struct {
	UserSrvI service.UserSrvInterface
}

// CheckEmail 检查邮箱
func (h *UserHandler) CheckEmail(c *gin.Context) {
	//初始化响应结构体
	entity := resp.EntityA{
		Code: 500,
		Msg:  "查询失败",
		Data: nil,
	}
	u := &model.User{Email: c.Param("email")}
	if u.Email == "" {
		entity.SetCodeAndMsg(500, "邮箱不能为空！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//查询数据库中是否已经有该邮箱
	r := h.UserSrvI.EmailExist(u)
	if r != nil {
		entity.SetCodeAndMsg(500, "该邮箱已被注册！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetCodeAndMsg(200, "该邮箱可用")
	c.JSON(200, gin.H{"entity": entity})
	return
}

// SendCode 发送验证码
func (h *UserHandler) SendCode(c *gin.Context) {
	//定义响应结构体
	entity := resp.EntityA{
		Code: 500,
		Msg:  "验证码发送失败",
		Data: nil,
	}
	//接收前端传来的email参数并放入Email字段
	Email := c.PostForm("email")
	if Email == "" {
		entity.SetCodeAndMsg(500, "邮箱不能为空")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//调用工具类生成随机数并转为string类型
	code := strconv.Itoa(utils.RandInt(123456, 999999))
	//为本次会话上下文初始化一个session
	session := sessions.Default(c)
	//设置session的有效期为五分钟，及其他属性
	session.Options(sessions.Options{MaxAge: 5 * 60, HttpOnly: true, Path: "/"})
	//将该验证码放入session中
	session.Set("code", code)
	session.Set("email", Email)
	//保存session
	if err := session.Save(); err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//调用service层发送验证码
	if err := h.UserSrvI.SendCode(code, Email); err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetCodeAndMsg(200, "验证码发送成功")
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *UserHandler) UserReg(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "用户注册失败",
		Data: nil,
	}
	// 若前端没有传回code则返回
	if c.Param("code") == "" {
		entity.SetCodeAndMsg(500, "验证码不能为空!")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	// 初始化一个user对象接收前端参数
	u := &model.User{}
	// 将前端传来的json序列化成user对象
	err := c.ShouldBindJSON(&u)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	addUser := &model.User{
		UserName: u.UserName,
		Password: u.Password,
		Email:    u.Email,
		Group:    "user",
	}
	// 创建session
	session := sessions.Default(c)
	// 取出该上下文中session中code的值
	code := session.Get("code")
	email := session.Get("email")
	//预防空指针
	if code == nil || email == nil {
		entity.Msg = "您的验证码已过期或Cookie状态异常！"
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	// 将参数中的code和session中取出的code作比较，其中code.(string)是将code（类型interface{}）断言为string类型
	if c.Param("code") != code.(string) || addUser.Email != email.(string) {
		entity.Msg = "验证码错误！"
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	// 验证成功后删除session中的code
	session.Delete("code")
	session.Delete("email")
	// 保存session
	sessionErr := session.Save()
	if sessionErr != nil {
		log.Println(sessionErr.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	// 调用service层添加该用户
	user, err := h.UserSrvI.AddUser(addUser)
	if err != nil {
		log.Println(err.Error())
		entity.SetCodeAndMsg(500, err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//设置响应信息
	entity.SetCodeAndMsg(200, "注册成功！")
	// 将user转换为响应user类型避免暴露过多字段
	entity.Data = user.ToRespUser()
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *UserHandler) UserLogin(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "登陆失败,请重试",
		Data: nil,
	}
	u := &model.User{}
	if err := c.ShouldBindJSON(u); err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	// 调用service层登录
	user, token, err := h.UserSrvI.UserLogin(u)
	if err != nil {
		log.Println(err.Error())
		entity.SetCodeAndMsg(500, err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	session := sessions.Default(c)
	session.Options(sessions.Options{MaxAge: 3600 * 24 * 14, HttpOnly: true, Path: "/"})
	// 将用户信息存入session
	session.Set("user", user)
	// 保存
	err = session.Save()
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	c.Header("Authorization", token)
	// 设置响应数据
	entity.SetCodeAndMsg(200, "登陆成功")
	entity.Token = token
	entity.Data = user.ToRespUser()
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
		entity.SetCodeAndMsg(500, "您的会话已过期，请重新登录!")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	u := &model.User{
		UserId: userInfo.(*model.User).UserId,
	}
	user, err := h.UserSrvI.UserInfo(u)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetCodeAndMsg(200, "OK")
	entity.SetEntityAndHeaderToken(c)
	entity.Data = user.ToRespUser()
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *UserHandler) Logout(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "Error",
		Data: nil,
	}
	session := sessions.Default(c)
	userInfo := session.Get("user")
	if userInfo == nil {
		entity.SetCodeAndMsg(500, "您的会话已过期，请重新登录!")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	fmt.Println(userInfo.(*model.User).UserId, "++++++++")
	h.UserSrvI.Logout(userInfo.(*model.User).UserId)
	session.Clear()
	if err := session.Save(); err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetCodeAndMsg(200, "退出成功！")
	c.Header("Authorization", "")
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *UserHandler) EditUserInfo(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "参数错误或Cookie状态异常！",
		Data: nil,
	}
	user := &model.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	updateUser := &model.User{
		UserId:   user.UserId,
		UserName: user.UserName,
		Icon:     user.Icon,
	}
	result := h.UserSrvI.EmailExist(user)
	if result == nil {
		entity.SetCodeAndMsg(500, "请检查参数")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//鉴定操作合法性
	if updateUser.UserId != result.UserId {
		entity.SetCodeAndMsg(500, "非法请求！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	if err := h.UserSrvI.EditUser(updateUser); err != nil {
		entity.SetCodeAndMsg(500, err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetCodeAndMsg(200, "更改信息成功!")
	entity.SetEntityAndHeaderToken(c)
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *UserHandler) EmailVerify(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "参数错误或Cookie状态异常！",
		Data: nil,
	}
	email := c.PostForm("email")
	code := c.PostForm("code")
	session := sessions.Default(c)
	if email == "" || code == "" || session.Get("email") == nil || session.Get("code") == nil {
		log.Println("参数错误或Cookie状态异常！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	if email != session.Get("email").(string) || code != session.Get("code").(string) {
		log.Println("验证码错误！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	session.Set("verifyEmail", email)
	session.Delete("code")
	session.Delete("email")
	if err := session.Save(); err != nil {
		log.Println(err.Error())
		entity.SetCodeAndMsg(500, "验证错误！请重试！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetEntityAndHeaderToken(c)
	entity.SetCodeAndMsg(200, "验证成功")
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *UserHandler) EditUserEmail(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "参数错误或Cookie状态异常！",
		Data: nil,
	}
	// 创建session
	session := sessions.Default(c)
	// 取出该上下文中session中code的值
	code := session.Get("code")
	oldEmail := session.Get("verifyEmail")
	newEmail := c.PostForm("newEmail")
	// 若前端没有传回code则返回
	if c.PostForm("code") == "" || c.PostForm("newEmail") == "" || code == nil || oldEmail == nil {
		log.Println("参数错误或Cookie状态异常！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	user := &model.User{}
	user.Email = oldEmail.(string)
	result := h.UserSrvI.EmailExist(user)
	if result == nil {
		log.Println("未查询到对应邮箱用户")
		entity.SetCodeAndMsg(500, "请检查邮箱后重试")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	// 将参数中的code和session中取出的code作比较，其中code.(string)是将code（类型interface{}）断言为string类型
	if c.PostForm("code") != code.(string) || newEmail != session.Get("email") {
		log.Println("验证码错误！")
		entity.Msg = "验证码错误！"
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	// 验证成功后删除session中的code
	session.Delete("code")
	session.Delete("email")
	session.Delete("oldEmail")
	// 保存session
	if sessionErr := session.Save(); sessionErr != nil {
		log.Println(sessionErr.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	updateUser := &model.User{
		UserId: result.UserId,
		Email:  newEmail,
	}
	if err := h.UserSrvI.EditUser(updateUser); err != nil {
		log.Println(err.Error())
		entity.SetCodeAndMsg(500, "修改邮箱失败,请重试")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetEntityAndHeaderToken(c)
	entity.SetCodeAndMsg(200, "更改邮箱成功!")
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *UserHandler) EditUserPw(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "参数错误或Cookie状态异常！",
		Data: nil,
	}
	// 创建session
	session := sessions.Default(c)
	// 取出该上下文中session中code的值
	code := session.Get("code")
	email := session.Get("email")
	// 若前端没有传回code则返回
	if c.PostForm("email") == "" || c.PostForm("code") == "" ||
		c.PostForm("password") == "" || code == nil || email == nil {
		log.Println("参数错误或Cookie状态异常！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	if c.PostForm("email") != email.(string) {
		entity.SetCodeAndMsg(500, "非法请求！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	// 将参数中的code和session中取出的code作比较，其中code.(string)是将code（类型interface{}）断言为string类型
	if c.PostForm("code") != code.(string) {
		entity.Msg = "验证码错误！"
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	// 验证成功后删除session中的code
	session.Delete("code")
	session.Delete("email")
	// 保存session
	if sessionErr := session.Save(); sessionErr != nil {
		log.Println(sessionErr.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	user := &model.User{}
	user.Email = c.PostForm("email")
	result := h.UserSrvI.EmailExist(user)
	if result == nil {
		log.Println("未查询到对应邮箱用户")
		entity.SetCodeAndMsg(500, "请检查参数")
		c.JSON(200, gin.H{"entity": entity})
		return
	}

	updateUser := &model.User{
		UserId:   result.UserId,
		Password: c.PostForm("password"),
	}
	if err := h.UserSrvI.EditUser(updateUser); err != nil {
		log.Println(err.Error())
		entity.SetCodeAndMsg(500, "修改密码失败,请重试")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetCodeAndMsg(200, "更改密码成功!")
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *UserHandler) AdminEditUser(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "参数错误！",
		Data: nil,
	}
	user := &model.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		log.Println("参数错误！")
		entity.SetCodeAndMsg(200, "参数错误！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	updateUser := &model.User{}
	updateUser.Email = user.Email
	result := h.UserSrvI.EmailExist(updateUser)
	if result == nil {
		entity.SetCodeAndMsg(500, "请检查参数")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//鉴权
	userGroup, _ := c.Get("userGroup")
	userId, _ := c.Get("userId")
	if userGroup.(string) != "Admin" || (userId != "Admin" && user.Group != result.Group) {
		log.Println("非法请求")
		entity.SetCodeAndMsg(500, "非法请求")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	if err := h.UserSrvI.EditUser(user); err != nil {
		log.Println("更新用户信息错误,请重试")
		entity.SetCodeAndMsg(500, "更新用户信息错误,请重试")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetCodeAndMsg(200, "更改用户信息成功！")
	entity.SetEntityAndHeaderToken(c)
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *UserHandler) EditUserStatus(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "参数错误！",
		Data: nil,
	}
	userId := c.PostForm("userId")
	status := c.PostForm("status")
	if userId == "" || status == "" {
		log.Println("参数错误,请重试")
		c.JSON(200, gin.H{"entity": entity})
	}
	//鉴权
	userGroup, _ := c.Get("userGroup")
	if userGroup.(string) != "Admin" || userId == "Admin" {
		log.Println("非法请求")
		entity.SetCodeAndMsg(500, "非法请求")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	if err := h.UserSrvI.EditUserStatus(userId, status); err != nil {
		entity.SetCodeAndMsg(500, "更改用户状态失败,请重试")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetCodeAndMsg(200, "更改用户信息成功")
	entity.SetEntityAndHeaderToken(c)
	c.JSON(200, gin.H{"entity": entity})
	return
}
