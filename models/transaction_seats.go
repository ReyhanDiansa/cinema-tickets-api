package models

type TransactionSeat struct {
	TransactionID string `json:"transaction_id"`
	SeatID        uint `json:"seat_id"`

	Seat     Seat     `gorm:"foreignKey:SeatID" json:"seat"`

}
