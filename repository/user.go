package repository

import (
	"P/model"
	"errors"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

type UserRepoInterface interface {
	EmailExist(user *model.User) *model.User
	Add(user *model.User) (*model.User, error)
	FindByEmail(user *model.User) (*model.User, error)
	FindByUserId(user *model.User) *model.User
	EditUserById(user *model.User) error
	EditUserStatusById(user *model.User) error
}

func (repo *UserRepository) FindByEmail(user *model.User) (*model.User, error) {
	result := repo.DB.Where("email=?", user.Email).Find(user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("未找到记录")
	}
	return user, nil
}

func (repo *UserRepository) FindByUserId(user *model.User) *model.User {
	var count int64 = 0
	repo.DB.Where("user_id=?", user.UserId).Find(user).Count(&count)
	if count > 0 {
		return user
	}
	return nil
}

func (repo *UserRepository) EmailExist(user *model.User) *model.User {
	var count int64 = 0
	repo.DB.Where("email = ?", user.Email).Find(user).Count(&count)
	if count > 0 {
		return user
	}
	return nil
}

func (repo *UserRepository) Add(user *model.User) (*model.User, error) {
	err := repo.DB.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) EditUserById(user *model.User) error {
	result := repo.DB.Model(user).Where("user_id=?", user.UserId).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *UserRepository) EditUserStatusById(user *model.User) error {
	if err := repo.DB.Model(&model.User{}).Where("user_id=?", user.UserId).Update("status", user.Status).Error; err != nil {
		return err
	}
	return nil
}
