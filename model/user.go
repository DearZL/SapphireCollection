package model

import (
	"P/resp"
	"gorm.io/gorm"
)

type User struct {
	UserId     string                `json:"userId" gorm:"column:user_id"`
	UserName   string                `json:"userName" gorm:"column:user_name"`
	Email      string                `json:"email" gorm:"column:email;unique"`
	Password   string                `json:"password" gorm:"column:password"`
	Icon       string                `json:"icon" gorm:"column:icon"`
	Status     bool                  `json:"status" gorm:"column:status;default:0"`
	Group      string                `json:"group" gorm:"default:user"`
	Collection []*resp.UserCommodity `json:"collection" gorm:"-"`
	gorm.Model
}

func (u *User) ToRespUser() *resp.User {
	re := &resp.User{
		UserId:     u.UserId,
		UserName:   u.UserName,
		Email:      u.Email,
		Icon:       u.Icon,
		Collection: u.Collection,
	}
	return re
}
