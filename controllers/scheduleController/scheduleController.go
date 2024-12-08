package schedulecontroller

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
	var schedules []models.Schedule
	var total int64
	pagination := utils.GetPagination(c)

	offset := (pagination.Page - 1) * pagination.Limit
	config.DB.Model(&models.Schedule{}).Count(&total)
	config.DB.Limit(pagination.Limit).Offset(offset).Preload("Cinema", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Name", "Location").Omit("CreatedAt", "UpdatedAt")
	}).Find(&schedules)

	if len(schedules) == 0 {
		utils.ResponseFormatter(c, 404, false, "no schedule data", nil)
		return
	}

	totalPages := (total + int64(pagination.Limit) - 1) / int64(pagination.Limit)
	paginatedResponse := utils.ResponsePaginationFormatter(schedules, pagination.Page, pagination.Limit, total, totalPages)

	utils.ResponseFormatter(c, 200, true, "success get schedules", paginatedResponse)
}

func Show(c *gin.Context) {
	var schedule models.Schedule

	id := c.Param("id")

	if err := config.DB.Preload("Film").Preload("Cinema").First(&schedule, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			utils.ResponseFormatter(c, 404, false, "no schedule data with that Id", nil)
			return

		default:
			utils.ResponseFormatter(c, 500, false, err.Error(), nil)
			return
		}
	}

	utils.ResponseFormatter(c, 200, true, "success get schedule", schedule)

}

func Create(c *gin.Context) {
	var schedule models.Schedule

	if err := c.ShouldBindJSON(&schedule); err != nil {
		utils.ResponseFormatter(c, 400, false, err.Error(), nil)
		return
	}

	var cinema models.Cinema
	if err := config.DB.First(&cinema, schedule.CinemaID).Error; err != nil {
		utils.ResponseFormatter(c, 400, false, "Cinema ID is invalid or does not exist", nil)
		return
	}
	schedule.Cinema = cinema

	// Validate FilmID
	var film models.Film
	if err := config.DB.First(&film, schedule.FilmID).Error; err != nil {
		utils.ResponseFormatter(c, 400, false, "Film ID is invalid or does not exist", nil)
		return
	}
	schedule.Film = film

	// Validate the schedule data
	err := validate.Struct(schedule)
	if err != nil {
		errors := utils.TranslateError(err, trans)
		utils.ResponseFormatter(c, 400, false, fmt.Sprintf("Validation error: %s", errors[0]), nil)
		return
	}

	config.DB.Create(&schedule)
	utils.ResponseFormatter(c, 201, true, "Success add schedule", schedule)
}

func Update(c *gin.Context) {
	var schedule models.Schedule
	id := c.Param("id")

	if err := c.ShouldBindJSON(&schedule); err != nil {
		utils.ResponseFormatter(c, 400, false, err.Error(), nil)
		return
	}

	var cinema models.Cinema
	if err := config.DB.First(&cinema, schedule.CinemaID).Error; err != nil {
		utils.ResponseFormatter(c, 400, false, "Cinema ID is invalid or does not exist", nil)
		return
	}
	schedule.Cinema = cinema

	var film models.Film
	if err := config.DB.First(&film, schedule.CinemaID).Error; err != nil {
		utils.ResponseFormatter(c, 400, false, "Cinema ID is invalid or does not exist", nil)
		return
	}
	schedule.Film = film

	// Validate cinema data
	err := validate.Struct(schedule)
	if err != nil {
		errors := utils.TranslateError(err, trans)
		utils.ResponseFormatter(c, 400, false, fmt.Sprintf("Validation error: %s", errors[0]), nil)
		return
	}

	if config.DB.Model(&schedule).Where("id = ?", id).Updates(&schedule).RowsAffected == 0 {
		utils.ResponseFormatter(c, 400, false, "cannot update schedule", nil)
		return
	}

	IdInt, _ := strconv.Atoi(id)
	schedule_id := uint(IdInt)

	schedule.ID = schedule_id
	utils.ResponseFormatter(c, 200, true, "success update schedule", schedule)
}

func Delete(c *gin.Context) {
	var schedule models.Schedule
	id := c.Param("id")

	if config.DB.Delete(&schedule, id).RowsAffected == 0 {
		utils.ResponseFormatter(c, 400, false, "cannot delete schedule", nil)
		return
	}

	utils.ResponseFormatter(c, 200, true, "success delete schedule", nil)

}
