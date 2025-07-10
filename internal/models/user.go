package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username        string `json:"username" gorm:"unique; not null"`
	Email           string `json:"email" gorm:"uniqueIndex; not null"`
	Password        string `json:"-" gorm:"not null"`
	IsEmailVerified bool   `json:"is_email_verified" gorm:"default:false"`
}
