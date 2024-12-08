package loginController

import (
	"cinema-tickets/config"
	"cinema-tickets/models"
	"cinema-tickets/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var validate = validator.New()

// validation error message translation config
var english = en.New()
var uni = ut.New(english, english)
var trans, _ = uni.GetTranslator("en")
var _ = enTranslations.RegisterDefaultTranslations(validate, trans)

func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the JSON body into loginData
	if err := c.ShouldBindJSON(&loginData); err != nil {
		utils.ResponseFormatter(c, 400, false, err.Error(), nil)
		return
	}

	// Find user by email
	var user models.User
	if err := config.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		utils.ResponseFormatter(c, 400, false, "Email is incorrect", nil)
		return
	}

	// Verify password
	if !utils.VerifyPassword(user.Password, loginData.Password) {
		utils.ResponseFormatter(c, 400, false, "Password is incorrect", nil)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.Email, user.Role, user.Name)
	if err != nil {
		utils.ResponseFormatter(c, 500, true, "Failed to generate token", nil)
		return
	}

	// Send back the JWT token
	utils.ResponseFormatter(c, 200, true, "Login Successfully", token)

}

func Create(c *gin.Context) {

	var user models.User
	user.Role = "user"
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
