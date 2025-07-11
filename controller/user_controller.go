package controller

import (
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

	return ctx.JSON(http.StatusCreated, response.UserResponse{Message: "User created", Data: user})
}

func (c *UserController) GetAllUsers(ctx echo.Context) error {
	var req request.PaginationRequest

	if err := ctx.Bind(&req); err != nil || req.Page <= 0 || req.Limit <= 0 {
		req.Page = 1
		req.Limit = 10
	}

	users, total, err := c.manager.GetAllUsers(req.Page, req.Limit, req.MinAge, req.NameContains)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.UserResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, response.UserResponse{
		Message: "Users fetched",
		Data: map[string]interface{}{
			"page":         req.Page,
			"limit":        req.Limit,
			"total":        total,
			"minAge":       req.MinAge,
			"nameContains": req.NameContains,
			"users":        users,
		},
	})
}

func (c *UserController) GetUserByID(ctx echo.Context) error {
	id := ctx.Param("id")
	user, err := c.manager.GetUserByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, response.UserResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, response.UserResponse{Message: "User fetched", Data: user})
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

	_, err := c.manager.GetUserByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, response.UserResponse{Message: "User already deleted"})
	}

	if err := c.manager.DeleteUser(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.UserResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, response.UserResponse{Message: "User soft-deleted"})
}
