package transactioncontroller

import (
	"cinema-tickets/config"
	"cinema-tickets/models"
	"cinema-tickets/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	var transactions []models.Transaction
	var total int64
	pagination := utils.GetPagination(c)

	offset := (pagination.Page - 1) * pagination.Limit
	config.DB.Model(&models.Transaction{}).Count(&total)
	config.DB.Limit(pagination.Limit).Offset(offset).Preload("User").Preload("Schedule").Preload("Schedule.Film").Preload("Schedule.Cinema").Preload("TransactionSeat.Seat").Find(&transactions)

	if len(transactions) == 0 {
		utils.ResponseFormatter(c, 404, false, "no cinema data", nil)
		return
	}

	totalPages := (total + int64(pagination.Limit) - 1) / int64(pagination.Limit)
	paginatedResponse := utils.ResponsePaginationFormatter(transactions, pagination.Page, pagination.Limit, total, totalPages)

	utils.ResponseFormatter(c, 200, true, "success get cinemas", paginatedResponse)
}

func CreateTransaction(c *gin.Context) {
	var input struct {
		UserID     uint   `json:"user_id" binding:"required"`
		ScheduleID uint   `json:"schedule_id" binding:"required"`
		SeatIDs    []uint `json:"seat_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ResponseFormatter(c, 400, false, err.Error(), nil)
		return
	}

	// Validate schedule
	var schedule models.Schedule
	if err := config.DB.Preload("Film").Preload("Cinema").First(&schedule, input.ScheduleID).Error; err != nil {
		utils.ResponseFormatter(c, 400, false, "Schedule not found", nil)
		return
	}

	//validate seat is in same cinema with schedule
	for _, seat := range input.SeatIDs {
		var seatToBook models.Seat
		if err := config.DB.First(&seatToBook, seat).Error; err != nil {
			utils.ResponseFormatter(c, 404, false, fmt.Sprintf("Seat with id %s not found", strconv.FormatUint(uint64(seat), 10)), nil)
			return
		}
		if seatToBook.CinemaID != schedule.CinemaID {
			utils.ResponseFormatter(c, 404, false, fmt.Sprintf("Seat with id %s is not in cinema with that schedule", strconv.FormatUint(uint64(seat), 10)), nil)
			return
		}

	}

	//validate id schedule has passed
	dateString := schedule.Time
	layout := "2006-01-02 15:04:05"
	givenTime, err := time.Parse(layout, dateString)
	if err != nil {
		utils.ResponseFormatter(c, 400, false, "Error Parsing Schedule time", nil)
		return
	}
	now := time.Now()
	if now.After(givenTime) {
		utils.ResponseFormatter(c, 400, false, "Cannot book ticket, The scheduled time has passed", nil)
		return
	}

	var alreadyBookedSeats []models.TransactionSeat

	// Check if any of the seats are already booked for the given schedule
	if err := config.DB.
		Table("transaction_seats").
		Select("transaction_seats.seat_id").
		Joins("JOIN transactions ON transaction_seats.transaction_id = transactions.transaction_id").
		Where("transaction_seats.seat_id IN ? AND transactions.schedule_id = ? AND transactions.status != ?", input.SeatIDs, schedule.ID, "canceled").
		Find(&alreadyBookedSeats).Error; err != nil {
		utils.ResponseFormatter(c, 500, false, "Error checking seat availability", nil)
		return
	}

	// If any seats are already booked, return an error response
	if len(alreadyBookedSeats) > 0 {
		for _, seat := range alreadyBookedSeats {
			utils.ResponseFormatter(c, 400, false, fmt.Sprintf("Seat with id %s is already booked", strconv.FormatUint(uint64(seat.SeatID), 10)), nil)
			return
		}
	}

	// Calculate total price
	totalPrice := int64(len(input.SeatIDs)) * schedule.Price

	// Create transaction
	transaction := models.Transaction{
		UserID:        input.UserID,
		ScheduleID:    input.ScheduleID,
		TotalPrice:    totalPrice,
		Status:        "pending",
		TransactionID: utils.GenerateTransactionID(),
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		utils.ResponseFormatter(c, 400, false, "Failed to create transaction", nil)
		return
	}

	// create data in transaction_seats
	for _, seat := range input.SeatIDs {
		var transaction_seat models.TransactionSeat
		transaction_seat.SeatID = seat
		transaction_seat.TransactionID = transaction.TransactionID

		if err := config.DB.Create(&transaction_seat).Error; err != nil {
			utils.ResponseFormatter(c, 400, false, fmt.Sprintf("Failed to add seat booking data on seat with id %s", strconv.FormatUint(uint64(seat), 10)), nil)
			return
		}
	}

	var TransactionSeat []models.TransactionSeat
	config.DB.Preload("Seat.Cinema").Where("transaction_id = ?", transaction.TransactionID).Find(&TransactionSeat)

	var user models.User
	if err := config.DB.First(&user, input.UserID).Error; err != nil {
		utils.ResponseFormatter(c, 400, false, "user not found", nil)
		return
	}
	

	transaction.Schedule = schedule
	transaction.User = user
	transaction.TransactionSeat = TransactionSeat

	utils.ResponseFormatter(c, 200, true, "Transaction created successfully", transaction)

}

func CancelTransaction(c *gin.Context) {
	transactionID := c.Param("id")

	// Find transaction
	var transaction models.Transaction
	if err := config.DB.First(&transaction, "transaction_id = ?", transactionID).Error; err != nil {
		utils.ResponseFormatter(c, 404, false, "Transaction not found", nil)
		return
	}

	if transaction.Status != "pending" {
		utils.ResponseFormatter(c, 400, false, "Only pending transactions can be canceled", nil)
		return
	}

	// Update transaction status
	transaction.Status = "canceled"
	if err := config.DB.Save(&transaction).Error; err != nil {
		utils.ResponseFormatter(c, 400, false, "Failed to cancel transaction", nil)
		return
	}

	utils.ResponseFormatter(c, 200, true, fmt.Sprintf("Transaction %s canceled successfully", transactionID), nil)

}

func CheckAvailableSeat(c *gin.Context) {
	var input struct {
		ScheduleID uint `json:"schedule_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ResponseFormatter(c, 400, false, err.Error(), nil)
		return
	}

	var availableSeats []models.Seat

	var schedule models.Schedule
	if err := config.DB.Where("id = ?", input.ScheduleID).First(&schedule).Error; err != nil {
		utils.ResponseFormatter(c, 400, false, "Schedule not found", nil)
		return
	}

	if err := config.DB.
		Table("seats").
		Select("seats.*").
		Where("seats.id NOT IN (?) AND seats.cinema_id = ?",
			config.DB.
				Table("transaction_seats").
				Select("transaction_seats.seat_id").
				Joins("JOIN transactions ON transaction_seats.transaction_id = transactions.transaction_id").
				Where("transactions.schedule_id = ? AND transactions.status != ?", input.ScheduleID, "canceled"), schedule.CinemaID,
		).Preload("Cinema").
		Find(&availableSeats).Error; err != nil {
		utils.ResponseFormatter(c, 500, false, "Error checking seat availability", nil)
		return
	}

	// If any seats are already booked, return an error response
	if len(availableSeats) > 0 {
		utils.ResponseFormatter(c, 200, true, "Success get available seat", availableSeats)
		return
	} else {
		utils.ResponseFormatter(c, 404, false, "No seat available", nil)
		return
	}
}

func Find(c *gin.Context) {
	var transactions models.Transaction
	transactionID := c.Param("transaction_id")

	if err := config.DB.Preload("User").Preload("Schedule").Preload("Schedule.Film").Preload("Schedule.Cinema").Preload("TransactionSeat.Seat").First(&transactions, "transaction_id = ? ",transactionID).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			utils.ResponseFormatter(c, 404, false, "no transaction data with that Id", nil)
			return

		default:
			utils.ResponseFormatter(c, 500, false, err.Error(), nil)
			return
		}
	}

	utils.ResponseFormatter(c, 200, true, "success get transaction", transactions)
}