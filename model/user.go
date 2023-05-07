package model

import (
	"P/resp"
	"gorm.io/gorm"
)

type User struct {
	UserId    string `json:"userId" gorm:"column:user_id"`
	UserName  string `json:"userName" gorm:"column:user_name"`
	Email     string `json:"email" gorm:"column:email" binding:"required"`
	Password  string `json:"password" gorm:"column:password"`
	Icon      string `json:"icon" gorm:"column:icon"`
	IsDeleted bool   `json:"isDeleted" gorm:"column:is_deleted;default:0"`
	Status    bool   `json:"status" gorm:"column:status;default:0"`
	gorm.Model
}

func (u *User) ToRespUser() *resp.User {
	re := &resp.User{
		UserId:   u.UserId,
		UserName: u.UserName,
		Email:    u.Email,
		Icon:     u.Icon,
	}
	return re
}
