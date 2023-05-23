package service

import (
	"P/model"
	"P/repository"
	"P/utils"
	"errors"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v5"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"log"
	"time"
)

type UserSrvInterface interface {
	EmailExist(user *model.User) *model.User
	Add(user *model.User) (*model.User, error)
	SendCode(code string, email string) error
	UserLogin(user *model.User) (*model.User, string, error)
	UserInfo(user *model.User) (*model.User, error)
	Logout(sid string) error
	EditUser(user *model.User) error
	EditUserStatus(userId string, status string) error
}

type UserService struct {
	Redis *redis.Client
	//userDao
	UserRepo repository.UserRepoInterface
	//sessionDao
	SessionRepo   repository.SessionRepositoryInterface
	CommodityRepo repository.CommodityRepoInterface
}

// EmailExist 判断邮箱是否已存在
func (srv *UserService) EmailExist(user *model.User) *model.User {
	return srv.UserRepo.EmailExist(user)
}

// Add 添加用户
func (srv *UserService) Add(user *model.User) (*model.User, error) {
	// 如果不含有email值则非法
	if user.Email == "" || user.Password == "" {
		return nil, errors.New("邮箱或密码不能为空！")
	}
	if result := srv.EmailExist(user); result != nil {
		return nil, errors.New("该邮箱已被注册！")
	}
	user.UserId = uuid.NewV4().String()
	user, err := srv.UserRepo.Add(user)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return user, errors.New("注册失败！")
}

// SendCode 发送验证码
func (srv *UserService) SendCode(code string, email string) error {
	err := utils.SendMail(email, "您的验证码为:", code+"     五分钟内有效!")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// UserLogin 用户登录
func (srv *UserService) UserLogin(user *model.User) (*model.User, string, error) {
	if user.Email == "" || user.Password == "" {
		return nil, "", errors.New("邮箱密码不能为空！")
	}
	// 向数据库查询该邮箱是否存在
	r := srv.EmailExist(user)
	if r == nil {
		return nil, "", errors.New("该用户不存在，请检查邮箱或注册！")
	}
	if r.Status != false {
		return nil, "", errors.New("该用户已被禁用")
	}
	u := &model.User{}
	u.Email = user.Email
	u, err := srv.UserRepo.FindByEmail(u)
	if err != nil {
		return nil, "", errors.New("登陆失败,请重试")
	}
	// 比对密码
	if u.Password != user.Password {
		return nil, "", errors.New("密码或邮箱错误，请检查后重新输入！")
	}
	var d time.Duration
	// 鉴别身份以生成不同有效期的token
	if user.UserId == "Admin" {
		d = time.Duration(viper.GetInt64("loginTimeout.admin")) * time.Minute
	} else {
		d = time.Duration(viper.GetInt64("loginTimeout.user")) * time.Hour
	}
	t := time.Now().Add(d)
	// 初始化一个claims结构体用以生成token
	claims := utils.UserClaims{
		UserId:    user.UserId,
		UserGroup: user.Group,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t),
		},
	}
	// 调用工具类生成token
	token, err := claims.GToken()
	if err != nil {
		return nil, "", errors.New("登陆失败,请重试")
	}
	// 向redis中写入userID和token
	srv.Redis.Set(user.UserId, token, d)
	return user, token, nil
}

// UserInfo 查询用户信息
func (srv *UserService) UserInfo(user *model.User) (*model.User, error) {
	u := srv.UserRepo.FindByUserId(user)
	if u == nil {
		return nil, errors.New("未找到用户,请检查参数")
	}
	comS, err := srv.CommodityRepo.FindComSByUserId(user.UserId)
	if err != nil {
		log.Println(err.Error())
	}
	commodities := model.Commodities{Commodities: comS}
	u.Collection = commodities.ToRespUserCommodities().UserCommodities
	return u, nil
}

// Logout 退出登录
func (srv *UserService) Logout(sid string) error {
	_, err := srv.SessionRepo.DeleteSession(sid)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// EditUser 编辑用户信息
func (srv *UserService) EditUser(user *model.User) error {
	if err := srv.UserRepo.EditUserById(user); err != nil {
		return err
	}
	return nil
}

func (srv *UserService) EditUserStatus(userId string, status string) error {
	user := &model.User{UserId: userId}
	result := srv.UserRepo.FindByUserId(user)
	if result == nil {
		return errors.New("未找到该用户")
	}
	if result.Group == "Admin" {
		return errors.New("无法对管理员操作")
	}
	if status == "1" {
		result.Status = true
	} else {
		result.Status = false
	}
	if err := srv.UserRepo.EditUserStatusById(result); err != nil {
		return err
	}
	return nil
}
