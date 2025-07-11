package main

import (
	"crud-app/config"
	"crud-app/controller"
	"crud-app/database"
	"crud-app/logger"
	"crud-app/manager"
	"crud-app/validation"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

var userController *controller.UserController

func init() {
	logger.InitLogger()
	config.InitializeConfig()
	database.Connect()
	validation.InitValidator()

	userManager := new(manager.UserManager)
	userController = controller.NewUserController(userManager)
}

func main() {
	e := echo.New()
	e.POST("/users", userController.CreateUser)
	e.GET("/users", userController.GetAllUsers)
	e.GET("/users/:id", userController.GetUserByID)
	e.PUT("/users/:id", userController.UpdateUser)
	e.DELETE("/users/:id", userController.DeleteUser)

	port := "8080"

	log.Printf("Server running at http://localhost:%s", port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
