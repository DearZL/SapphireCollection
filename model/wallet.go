package model

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	UserId  string
	Balance float64
}
