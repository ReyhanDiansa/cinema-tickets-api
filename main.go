package main

import (
	"cinema-tickets/config"
	authcontroller "cinema-tickets/controllers/authController"
	cinemacontroller "cinema-tickets/controllers/cinemaController"
	filmcontroller "cinema-tickets/controllers/filmController"
	schedulecontroller "cinema-tickets/controllers/scheduleController"
	seatcontroller "cinema-tickets/controllers/seatController"
	transactioncontroller "cinema-tickets/controllers/transactionController"
	usercontroller "cinema-tickets/controllers/userController"
	middleware "cinema-tickets/middlewares"

	// "cinema-tickets/controllers/transactionController"
	"cinema-tickets/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config.ConnectDB()
	models.AutoMigrate()

	app := gin.Default()

	route := app.Group("/api")
	{
		auth := route.Group("/auth")
		{
			auth.POST("/login", authcontroller.Login)
			auth.POST("/register", authcontroller.Create)
		}
		//film endpoint
		film := route.Group("/film", middleware.AuthMiddleware())
		{
			film.GET("/", filmcontroller.Index)
			film.POST("/", middleware.RoleMiddleware("admin"), filmcontroller.Create)
			film.GET("/:id", filmcontroller.Show)
			film.PUT("/:id", middleware.RoleMiddleware("admin"), filmcontroller.Update)
			film.DELETE("/:id", middleware.RoleMiddleware("admin"), filmcontroller.Delete)
		}

		cinema := route.Group("/cinema", middleware.AuthMiddleware())
		{
			cinema.GET("/", cinemacontroller.Index)
			cinema.POST("/", middleware.RoleMiddleware("admin"), cinemacontroller.Create)
			cinema.GET("/:id", cinemacontroller.Show)
			cinema.PUT("/:id", middleware.RoleMiddleware("admin"), cinemacontroller.Update)
			cinema.DELETE("/:id", middleware.RoleMiddleware("admin"), cinemacontroller.Delete)
		}

		user := route.Group("/user", middleware.AuthMiddleware())
		{
			user.GET("/", middleware.RoleMiddleware("admin"), usercontroller.Index)
			user.POST("/", middleware.RoleMiddleware("admin"), usercontroller.Create)
			user.GET("/:id", middleware.RoleMiddleware("admin"), usercontroller.Show)
			user.PUT("/:id", middleware.RoleMiddleware("admin"), usercontroller.Update)
			user.DELETE("/:id", middleware.RoleMiddleware("admin"), usercontroller.Delete)
		}

		seat := route.Group("/seat", middleware.AuthMiddleware())
		{
			seat.GET("/", seatcontroller.Index)
			seat.POST("/", middleware.RoleMiddleware("admin"), seatcontroller.Create)
			seat.GET("/:id", seatcontroller.Show)
			seat.PUT("/:id", middleware.RoleMiddleware("admin"), seatcontroller.Update)
			seat.DELETE("/:id",middleware.RoleMiddleware("admin"),  seatcontroller.Delete)
		}

		schedule := route.Group("/schedule", middleware.AuthMiddleware())
		{
			schedule.GET("/", schedulecontroller.Index)
			schedule.POST("/", middleware.RoleMiddleware("admin"), schedulecontroller.Create)
			schedule.GET("/:id", schedulecontroller.Show)
			schedule.PUT("/:id", middleware.RoleMiddleware("admin"), schedulecontroller.Update)
			schedule.DELETE("/:id", middleware.RoleMiddleware("admin"), schedulecontroller.Delete)
		}

		transaction := route.Group("/transaction", middleware.AuthMiddleware())
		{
			transaction.GET("/", middleware.RoleMiddleware("admin"), transactioncontroller.Index)
			transaction.GET("/:transaction_id", transactioncontroller.Find)
			transaction.POST("/check-available", transactioncontroller.CheckAvailableSeat)
			transaction.POST("/book-ticket", transactioncontroller.CreateTransaction)
			transaction.POST("/cancel/:id", middleware.RoleMiddleware("admin"), transactioncontroller.CancelTransaction)
		}
	}

	app.Run(":8000")
}
