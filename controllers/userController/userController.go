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
	var users []models.User
	var total int64
	pagination := utils.GetPagination(c)

	offset := (pagination.Page - 1) * pagination.Limit
	config.DB.Model(&models.User{}).Count(&total)
	config.DB.Limit(pagination.Limit).Offset(offset).Preload("Cinema").Find(&users)

	if len(users) == 0 {
		utils.ResponseFormatter(c, 404, false, "no user data", nil)
		return
	}

	totalPages := (total + int64(pagination.Limit) - 1) / int64(pagination.Limit)
	paginatedResponse := utils.ResponsePaginationFormatter(users, pagination.Page, pagination.Limit, total,totalPages)

	utils.ResponseFormatter(c, 200, true, "success get users", paginatedResponse)
}

func Show(c *gin.Context) {
	var user models.User

	id := c.Param("id")

	if err := config.DB.First(&user, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			utils.ResponseFormatter(c, 404, false, "no user data with that Id", nil)
			return

		default:
			utils.ResponseFormatter(c, 500, false, err.Error(), nil)
			return
		}
	}

	utils.ResponseFormatter(c, 200, true, "success get user", user)
}

func Create(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ResponseFormatter(c, 400, false, err.Error(), nil)
		return
	}

	// Validate user data
	err := validate.Struct(user)
	if err != nil {
		errors := utils.TranslateError(err, trans)
		utils.ResponseFormatter(c, 400, false, fmt.Sprintf("Validation error: %s", errors), nil)
		return
	}

	result := config.DB.Where("email = ?", user.Email).First(&user)

	if result.Error == nil {
		utils.ResponseFormatter(c, 400, false, fmt.Sprintf("Validation error: User with Email %s already exists", user.Email), nil)
		return
	}

	user.Password = utils.HashPassword(user.Password)

	config.DB.Create(&user)
	utils.ResponseFormatter(c, 201, true, "success add user", user)
}

func Update(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ResponseFormatter(c, 400, false, err.Error(), nil)
		return
	}

	result := config.DB.Where("email = ?, id != ?", user.Email, id).First(&user)

	if result.Error == nil {
		utils.ResponseFormatter(c, 400, false, fmt.Sprintf("Validation error: User with Email %s already exists", user.Email), nil)
		return
	}

	user.Password = utils.HashPassword(user.Password)

	if config.DB.Model(&user).Where("id = ?", id).Updates(&user).RowsAffected == 0 {
		utils.ResponseFormatter(c, 400, false, "cannot update user", nil)
		return
	}

	IdInt, _ := strconv.Atoi(id)
	user_id := uint(IdInt)

	user.ID = user_id
	utils.ResponseFormatter(c, 200, true, "success update user", user)
}

func Delete(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if config.DB.Delete(&user, id).RowsAffected == 0 {
		utils.ResponseFormatter(c, 400, false, "cannot delete user", nil)
		return
	}

	utils.ResponseFormatter(c, 200, true, "success delete user", nil)

}
