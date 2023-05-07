package repository

import (
	"P/model"
	"gorm.io/gorm"
)

type SessionRepositoryInterface interface {
	DeleteSession(SId string) (string, error)
}

type SessionRepository struct {
	DB *gorm.DB
}

func (repo SessionRepository) DeleteSession(SId string) (string, error) {
	m := model.Sessions{}
	err := repo.DB.Where("id=?", SId).Delete(&m).Error
	if err != nil {
		return "", err
	}
	return "OK", nil
}
