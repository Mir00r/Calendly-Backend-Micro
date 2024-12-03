package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

// TableName overrides the default table name
func (User) TableName() string {
	return "auth.users"
}
