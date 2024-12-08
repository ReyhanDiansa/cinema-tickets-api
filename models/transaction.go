package models

import "time"

type Transaction struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null" json:"user_id"`
	ScheduleID     uint      `gorm:"not null" json:"schedule_id"`
	TotalPrice     int64   `gorm:"not null" json:"total_price"`
	Status         string    `gorm:"type:enum('pending', 'completed', 'canceled');default:'pending'" json:"status"`
	TransactionID  string    `gorm:"type:varchar(255);unique;not null" json:"transaction_id"`
	TransactionSeat  []TransactionSeat `gorm:"foreignKey:TransactionID;references:TransactionID" json:"transaction_seats"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	User     User     `gorm:"foreignKey:UserID" json:"user"`
	Schedule Schedule `gorm:"foreignKey:ScheduleID" json:"schedule"`
}
