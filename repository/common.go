package repository

import (
	"P/model"
	"errors"
	"gorm.io/gorm"
	"log"
)

type CommonRepository struct {
	DB *gorm.DB
}
type CommonRepositoryInterface interface {
	AddFile(files []model.File) error
	FindFile(filename string) (*model.File, error)
}

func (repo *CommonRepository) AddFile(files []model.File) error {
	err := repo.DB.Create(&files).Error
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
func (repo *CommonRepository) FindFile(filename string) (*model.File, error) {
	f := &model.File{}
	result := repo.DB.Where("filename=?", filename).Find(&f)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("未找到记录")
	}
	return f, nil
}
