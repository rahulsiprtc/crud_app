package controller

import (
	"fmt"
	"net/http"

	"crud-app/manager"
	"crud-app/request"
	"crud-app/response"
	"crud-app/validation"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	manager *manager.UserManager
}

func NewUserController(m *manager.UserManager) *UserController {
	return &UserController{
		manager: m,
	}
}

func (c *UserController) CreateUser(ctx echo.Context) error {
	var req request.CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.UserResponse{Message: "Invalid request"})
	}

	if err := validation.Validator.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.UserResponse{Message: err.Error()})
	}
	user, err := c.manager.CreateUser(req)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.UserResponse{Message: err.Error()})
	}
	fmt.Print(user)

	// return ctx.JSON(http.StatusOK, user, response.UserResponse{Message:"user created" })
	return ctx.JSON(http.StatusCreated, response.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Age:     user.Age,
		Message: "User created",
	})
}

// done
func (c *UserController) GetAllUsers(ctx echo.Context) error {
	var req request.PaginationRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.UserResponse{Message: "Invalid request", Error: err.Error()})
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// result, err := c.manager.GetAllUsers(req.ID, req.Page, req.Limit, req.MinAge, req.NameContains)
	result, err := c.manager.GetAllUsers(req)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.UserResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, result)
}

func (c *UserController) UpdateUser(ctx echo.Context) error {
	id := ctx.Param("id")
	var req request.UpdateUserRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.UserResponse{Message: "Invalid request"})
	}

	if err := validation.Validator.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.UserResponse{Message: err.Error()})
	}

	msg, err := c.manager.UpdateUser(id, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.UserResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, response.UserResponse{Message: msg})

}

func (c *UserController) DeleteUser(ctx echo.Context) error {
	id := ctx.Param("id")

	if err := c.manager.DeleteUser(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.UserResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, response.UserResponse{Message: "User deleted"})
}
