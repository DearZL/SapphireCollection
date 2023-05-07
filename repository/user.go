package repository

import (
	"P/model"
	"gorm.io/gorm"
	"log"
)

type UserRepository struct {
	DB *gorm.DB
}

type UserRepoInterface interface {
	EmailExist(user model.User) *model.User
	Add(user model.User) (*model.User, error)
	FindByEmail(user model.User) (*model.User, error)
	FindById(user model.User) (*model.User, error)
	EditUser(user model.User) error
}

func (repo *UserRepository) FindByEmail(user model.User) (*model.User, error) {
	err := repo.DB.Where("email=?", user.Email).Find(&user).Error
	if err != nil {
		log.Println("查找用户出现错误")
		return &model.User{}, err
	}
	return &user, nil
}

func (repo *UserRepository) FindById(user model.User) (*model.User, error) {
	err := repo.DB.Where("user_id=?", user.UserId).Find(&user).Error
	if err != nil {
		log.Println("查找用户出现错误")
		return &model.User{}, err
	}
	return &user, nil
}

func (repo *UserRepository) EmailExist(user model.User) *model.User {
	var count int64 = 0
	repo.DB.Where("email = ?", user.Email).Find(&user).Count(&count)
	if count > 0 {
		return &user
	}
	return nil
}

func (repo *UserRepository) Add(user model.User) (*model.User, error) {
	err := repo.DB.Create(&user).Error
	if err != nil {
		log.Println("用户注册失败")
		return nil, err
	}
	return &user, nil
}
func (repo *UserRepository) EditUser(user model.User) error {
	err := repo.DB.Model(&user).Where("user_id=?", user.UserId).Updates(&user).Error
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
