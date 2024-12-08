package models

import "time"

type Film struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title" validate:"required"`
	Description string    `gorm:"type:text" json:"description" validate:"required"`
	Duration    int       `gorm:"not null" json:"duration" validate:"required"`
	Genre       string    `gorm:"type:varchar(255);not null" json:"genre" validate:"required"`
	Rating      string    `gorm:"type:varchar(5);not null" json:"rating" validate:"required"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
