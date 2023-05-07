package service

import (
	"P/model"
	"P/repository"
	"P/utils"
	"errors"
	"github.com/satori/go.uuid"
	"log"
)

type UserSrvInterface interface {
	EmailExist(user model.User) *model.User
	Add(user model.User) (*model.User, error)
	SendCode(code string, user model.User) error
	UserLogin(user model.User) (*model.User, error)
	UserInfo(user model.User) (*model.User, error)
	Logout(sid string) error
	EditUser(user model.User) error
}

type UserService struct {
	//userDao
	UserRepo repository.UserRepoInterface
	//sessionDao
	SessionRepo repository.SessionRepositoryInterface
}

// EmailExist 判断邮箱是否已存在
func (srv *UserService) EmailExist(user model.User) *model.User {
	return srv.UserRepo.EmailExist(user)
}

// Add 添加用户
func (srv *UserService) Add(user model.User) (*model.User, error) {
	if result := srv.EmailExist(user); result != nil {
		return nil, errors.New("该邮箱已被注册！")
	}
	user.UserId = uuid.NewV4().String()
	if user.UserId == "" {
		return nil, errors.New("id生成错误")
	}
	return srv.UserRepo.Add(user)
}

// SendCode 发送验证码
func (srv *UserService) SendCode(code string, user model.User) error {
	err := utils.SendMail(user.Email, "您的验证码为:", code+"     五分钟内有效!")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// UserLogin 用户登录
func (srv *UserService) UserLogin(user model.User) (*model.User, error) {
	return srv.UserRepo.FindByEmail(user)
}

// UserInfo 查询用户信息
func (srv *UserService) UserInfo(user model.User) (*model.User, error) {
	u, err := srv.UserRepo.FindById(user)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
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
func (srv *UserService) EditUser(user model.User) error {
	return srv.UserRepo.EditUser(user)
}
