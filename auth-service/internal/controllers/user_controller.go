package controllers

import (
	"net/http"

	"github.com/MarkoLuna/AuthService/internal/models"
	"github.com/MarkoLuna/AuthService/internal/services"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type UserController struct {
	userService services.UserService
	validator   *validator.Validate
}

func NewUserController(userService services.UserService) UserController {
	return UserController{
		userService: userService,
		validator:   validator.New(),
	}
}

// CreateUser
// @Tags Users
// @Summary Create a new user
// @Description Create a new user with username and password
// @Accept json
// @Produce json
// @Param user body models.UserRequest true "User data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} error
// @Failure 409 {object} error
// @Security ApiKeyAuth
// @Router /api/user/ [post]
func (ctrl UserController) CreateUser(c echo.Context) error {
	var req models.UserRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	if err := ctrl.validator.Struct(req); err != nil {
		return c.String(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	user, err := ctrl.userService.CreateUser(req)
	if err != nil {
		if err.Error() == "username already exists" {
			return c.String(http.StatusConflict, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}

	response := user.ToResponse()
	return c.JSON(http.StatusCreated, response)
}

// GetUsers
// @Tags Users
// @Summary List all users
// @Description Get a list of all registered users
// @Produce json
// @Success 200 {array} models.UserResponse
// @Security ApiKeyAuth
// @Router /api/user/ [get]
func (ctrl UserController) GetUsers(c echo.Context) error {
	users, err := ctrl.userService.GetUsers()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

// GetUserById
// @Tags Users
// @Summary Get user by ID
// @Description Get a single user by their ID
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 404 {object} error
// @Security ApiKeyAuth
// @Router /api/user/{userId} [get]
func (ctrl UserController) GetUserById(c echo.Context) error {
	userId := c.Param("userId")
	user, err := ctrl.userService.GetUserById(userId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if user == nil {
		return c.String(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, user)
}

// UpdateUser
// @Tags Users
// @Summary Update a user
// @Description Update user details
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param user body models.UserRequest true "Updated user data"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Security ApiKeyAuth
// @Router /api/user/{userId} [put]
func (ctrl UserController) UpdateUser(c echo.Context) error {
	userId := c.Param("userId")

	var req models.UserRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	if err := ctrl.validator.Struct(req); err != nil {
		return c.String(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	user, err := ctrl.userService.UpdateUser(userId, req)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if user == nil {
		return c.String(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteUser
// @Tags Users
// @Summary Delete a user
// @Description Soft-delete (disable) a user by ID
// @Produce json
// @Param userId path string true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} error
// @Security ApiKeyAuth
// @Router /api/user/{userId} [delete]
func (ctrl UserController) DeleteUser(c echo.Context) error {
	userId := c.Param("userId")
	if err := ctrl.userService.DeleteUser(userId); err != nil {
		if err.Error() == "user not found" {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
