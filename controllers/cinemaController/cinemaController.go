package cinemacontroller

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
	var cinemas []models.Cinema
	var total int64
	pagination := utils.GetPagination(c)

	offset := (pagination.Page - 1) * pagination.Limit
	config.DB.Model(&models.Cinema{}).Count(&total)
	config.DB.Limit(pagination.Limit).Offset(offset).Find(&cinemas)

	if len(cinemas) == 0 {
		utils.ResponseFormatter(c, 404, false, "no cinema data", nil)
		return
	}

	totalPages := (total + int64(pagination.Limit) - 1) / int64(pagination.Limit)

	paginatedResponse := utils.ResponsePaginationFormatter(cinemas, pagination.Page, pagination.Limit, total,totalPages)

	utils.ResponseFormatter(c, 200, true, "success get cinemas", paginatedResponse)
}

func Show(c *gin.Context) {
	var cinema models.Cinema

	id := c.Param("id")

	if err := config.DB.First(&cinema, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			utils.ResponseFormatter(c, 404, false, "no cinema data with that Id", nil)
			return

		default:
			utils.ResponseFormatter(c, 500, false, err.Error(), nil)
			return
		}
	}

	utils.ResponseFormatter(c, 200, true, "success get cinema", cinema)

}

func Create(c *gin.Context){
	var cinema models.Cinema
	
	if err := c.ShouldBindJSON(&cinema); err != nil {
		utils.ResponseFormatter(c, 400, false, err.Error(), nil)
		return
	}

	// Validate cinema data
	err := validate.Struct(cinema)
	if err != nil {
		errors := utils.TranslateError(err, trans)
		utils.ResponseFormatter(c, 400, false, fmt.Sprintf("Validation error: %s", errors[0]), nil)
		return
	}

	config.DB.Create(&cinema)
	utils.ResponseFormatter(c, 201, true, "success add cinema", cinema)
}


func Update(c *gin.Context){
	var cinema models.Cinema
	id := c.Param("id")

	
	if err := c.ShouldBindJSON(&cinema); err != nil {
		utils.ResponseFormatter(c, 400, false, err.Error(), nil)
		return
	}

	if config.DB.Model(&cinema).Where("id = ?", id).Updates(&cinema).RowsAffected == 0 {
		utils.ResponseFormatter(c, 400, false, "cannot update cinema", nil)
		return
	}

	// Validate cinema data
	err := validate.Struct(cinema)
	if err != nil {
		errors := utils.TranslateError(err, trans)
		utils.ResponseFormatter(c, 400, false, fmt.Sprintf("Validation error: %s", errors[0]), nil)
		return
	}

	IdInt, _ := strconv.Atoi(id)
	user_id := uint(IdInt)

	cinema.ID = user_id
	
	utils.ResponseFormatter(c, 200, true, "success update cinema", cinema)
}

func Delete(c *gin.Context){
	var cinema models.Cinema
	id := c.Param("id")

	if config.DB.Delete(&cinema, id).RowsAffected == 0 {
		utils.ResponseFormatter(c, 400, false, "cannot delete cinema", nil)
		return
	}

	utils.ResponseFormatter(c, 200, true, "success delete cinema", nil)

}