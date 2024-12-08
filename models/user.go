package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name" validate:"required"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email" validate:"required,email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password" validate:"required,min=8"`
	Role      string    `gorm:"type:enum('admin', 'user');default:'user'" json:"role" validate:"required,oneof=admin user"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
