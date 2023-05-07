package repository

import (
	"P/model"
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
	err := repo.DB.Where("filename=?", filename).Find(&f).Error
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return f, nil
}
