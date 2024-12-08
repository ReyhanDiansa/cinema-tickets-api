package filmcontroller

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
	var films []models.Film
	var total int64
	pagination := utils.GetPagination(c)

	offset := (pagination.Page - 1) * pagination.Limit
	config.DB.Model(&models.Film{}).Count(&total)
	config.DB.Limit(pagination.Limit).Offset(offset).Preload("Cinema").Find(&films)

	if len(films) == 0 {
		utils.ResponseFormatter(c, 404, false, "no film data", nil)
		return
	}
	totalPages := (total + int64(pagination.Limit) - 1) / int64(pagination.Limit)

	paginatedResponse := utils.ResponsePaginationFormatter(films, pagination.Page, pagination.Limit, total,totalPages)
		

	utils.ResponseFormatter(c, 200, true, "success get films", paginatedResponse)
}

func Show(c *gin.Context) {
	var film models.Film

	id := c.Param("id")

	if err := config.DB.First(&film, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			utils.ResponseFormatter(c, 404, false, "no film data with that Id", nil)
			return

		default:
			utils.ResponseFormatter(c, 500, false, err.Error(), nil)
			return
		}
	}

	utils.ResponseFormatter(c, 200, true, "success get film", film)

}

func Create(c *gin.Context){
	var film models.Film
	
	if err := c.ShouldBindJSON(&film); err != nil {
		utils.ResponseFormatter(c, 400, false, err.Error(), nil)
		return
	}

	// Validate cinema data
	err := validate.Struct(film)
	if err != nil {
		errors := utils.TranslateError(err, trans)
		utils.ResponseFormatter(c, 400, false, fmt.Sprintf("Validation error: %s", errors[0]), nil)
		return
	}

	config.DB.Create(&film)
	utils.ResponseFormatter(c, 201, true, "success add film", film)
}


func Update(c *gin.Context){
	var film models.Film
	id := c.Param("id")

	
	if err := c.ShouldBindJSON(&film); err != nil {
		utils.ResponseFormatter(c, 400, false, err.Error(), nil)
		return
	}

	if config.DB.Model(&film).Where("id = ?", id).Updates(&film).RowsAffected == 0 {
		utils.ResponseFormatter(c, 400, false, "cannot update film", nil)
		return
	}

	// Validate cinema data
	err := validate.Struct(film)
	if err != nil {
		errors := utils.TranslateError(err, trans)
		utils.ResponseFormatter(c, 400, false, fmt.Sprintf("Validation error: %s", errors[0]), nil)
		return
	}

	IdInt, _ := strconv.Atoi(id)
	user_id := uint(IdInt)

	film.ID = user_id
	utils.ResponseFormatter(c, 200, true, "success update film", film)
}

func Delete(c *gin.Context){
	var film models.Film
	id := c.Param("id")

	if config.DB.Delete(&film, id).RowsAffected == 0 {
		utils.ResponseFormatter(c, 400, false, "cannot delete film", nil)
		return
	}

	utils.ResponseFormatter(c, 200, true, "success delete film", nil)

}