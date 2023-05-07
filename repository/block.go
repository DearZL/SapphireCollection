package repository

import "gorm.io/gorm"

type BlockRepoInterface interface {
}
type BlockRepository struct {
	DB *gorm.DB
}
