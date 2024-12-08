package models


type Seat struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Number   string `gorm:"type:varchar(10);not null" json:"number" validate:"required"`
	CinemaID uint   `gorm:"not null" json:"cinema_id" validate:"required"`

	Cinema Cinema `gorm:"foreignKey:CinemaID"`
}
