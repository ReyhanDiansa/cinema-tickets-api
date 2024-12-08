package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)


func GenerateTransactionID() string {
	rand.Seed(time.Now().UnixNano())
	return time.Now().Format("20060102") + "-" + RandString(8)
}

func RandString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func ResponseFormatter(c *gin.Context, code int, success bool, message string, data any) {
	c.JSON(code, gin.H{"success": success, "message": message, "data": data})
}

func HashPassword(password string) string {
	hash := md5.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func VerifyPassword(storedHash, inputPassword string) bool {
	inputHash := HashPassword(inputPassword)

	return storedHash == inputHash
}

func TranslateError(err error, trans ut.Translator) (errs []error) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errs = append(errs, translatedErr)
	}
	return errs
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func GetPagination(c *gin.Context) Pagination {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	return Pagination{Page: page, Limit: limit}
}

func ResponsePaginationFormatter(items any, page int, limit int, total int64, totalPage int64) gin.H {
    return gin.H{
        "items": items,
        "meta": gin.H{
            "current_page":  page,
            "per_page": limit,
            "total_data": total,
            "total_page": totalPage,
        },
    }
}

type Claims struct {
	Email string `json:"email"`
	Role string `json:"role"`
	Name string `json:"name"`
	jwt.StandardClaims
}
// Function to generate JWT
func GenerateJWT(email, role, name string) (string, error) {
var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 1 day
	claims := &Claims{
		Email: email,
		Role: role,
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
