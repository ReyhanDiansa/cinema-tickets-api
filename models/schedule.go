package models

import "time"

type Schedule struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	FilmID   uint    `gorm:"not null" json:"film_id"`
	CinemaID uint    `gorm:"not null" json:"cinema_id"`
	Time     string  `gorm:"not null;type:varchar(250)" json:"time"`
	Price    int64 `gorm:"not null" json:"price"`

	Film      Film      `gorm:"foreignKey:FilmID"`
	Cinema    Cinema    `gorm:"foreignKey:CinemaID"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
