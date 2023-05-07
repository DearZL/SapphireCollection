package model

import (
	"gorm.io/gorm"
	"time"
)

type Sessions struct {
	Id        string
	Data      string
	ExpiresAt time.Time
	gorm.Model
}
