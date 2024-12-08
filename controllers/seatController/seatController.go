package seatcontroller

import (
	"cinema-tickets/config"
	"cinema-tickets/models"
	"cinema-tickets/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"gorm.io/gorm"
)

var validate = validator.New()

// validation error message translation config
var english = en.New()
var uni = ut.New(english, english)
var trans, _ = uni.GetTranslator("en")
var _ = enTranslations.RegisterDefaultTranslations(validate, trans)

func Index(c *gin.Context) {
	var seats []models.Seat
	var total int64
	pagination := utils.GetPagination(c)


	offset := (pagination.Page - 1) * pagination.Limit
	config.DB.Model(&models.Seat{}).Count(&total)
	config.DB.Limit(pagination.Limit).Offset(offset).Preload("Cinema").Find(&seats)

	if len(seats) == 0 {
		utils.ResponseFormatter(c, 404, false, "no seat data", nil)
		return
	}

    totalPages := (total + int64(pagination.Limit) - 1) / int64(pagination.Limit)
	paginatedResponse := utils.ResponsePaginationFormatter(seats, pagination.Page, pagination.Limit, total,totalPages)

	utils.ResponseFormatter(c, 200, true, "success get seats", paginatedResponse)
}

func Show(c *gin.Context) {
	var seat models.Seat

	id := c.Param("id")

	if err := config.DB.Preload("Cinema").First(&seat, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			utils.ResponseFormatter(c, 404, false, "no seat data with that Id", nil)
			return

		default:
			utils.ResponseFormatter(c, 500, false, err.Error(), nil)
			return
		}
	}

	utils.ResponseFormatter(c, 200, true, "success get seat", seat)

}

func Create(c *gin.Context) {
	var seat models.Seat

	// Bind the input JSON to the Seat struct
	if err := c.ShouldBindJSON(&seat); err != nil {
		utils.ResponseFormatter(c, 400, false, err.Error(), nil)
		return
	}

	var cinema models.Cinema
	if err := config.DB.First(&cinema, seat.CinemaID).Error; err != nil {
		utils.ResponseFormatter(c, 400, false, "Cinema ID is invalid or does not exist", nil)
		return
	}

	seat.Cinema = cinema

	err := validate.Struct(seat)
	if err != nil {
		errors := utils.TranslateError(err, trans)
		utils.ResponseFormatter(c, 400, false, fmt.Sprintf("Validation error: %s", errors[0]), nil)
		return
	}


	config.DB.Create(&seat)
	utils.ResponseFormatter(c, 201, true, "success add seat", seat)
}

func Update(c *gin.Context) {
	var seat models.Seat
	id := c.Param("id")

	if err := c.ShouldBindJSON(&seat); err != nil {
		utils.ResponseFormatter(c, 400, false, err.Error(), nil)
		return
	}

	var cinema models.Cinema
	if err := config.DB.First(&cinema, seat.CinemaID).Error; err != nil {
		utils.ResponseFormatter(c, 400, false, "Cinema ID is invalid or does not exist", nil)
		return
	}

	seat.Cinema = cinema

	// Validate cinema data
	err := validate.Struct(seat)
	if err != nil {
		errors := utils.TranslateError(err, trans)
		utils.ResponseFormatter(c, 400, false, fmt.Sprintf("Validation error: %s", errors[0]), nil)
		return
	}
	
	if config.DB.Model(&seat).Where("id = ?", id).Updates(&seat).RowsAffected == 0 {
		utils.ResponseFormatter(c, 400, false, "cannot update seat", nil)
		return
	}


	IdInt, _ := strconv.Atoi(id)
	user_id := uint(IdInt)

	seat.ID = user_id
	utils.ResponseFormatter(c, 200, true, "success update seat", seat)
}

func Delete(c *gin.Context) {
	var seat models.Seat
	id := c.Param("id")

	if config.DB.Delete(&seat, id).RowsAffected == 0 {
		utils.ResponseFormatter(c, 400, false, "cannot delete seat", nil)
		return
	}

	utils.ResponseFormatter(c, 200, true, "success delete seat", nil)

}
