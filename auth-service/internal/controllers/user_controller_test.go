package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MarkoLuna/AuthService/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type userServiceSuccessStub struct {
	users map[string]models.User
}

func newUserServiceSuccessStub() *userServiceSuccessStub {
	u := map[string]models.User{
		"user-1": {
			Id:        "user-1",
			Username:  "existing",
			Email:     "existing@test.com",
			FirstName: "Existing",
			LastName:  "User",
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	return &userServiceSuccessStub{users: u}
}

func (s *userServiceSuccessStub) GetUserId(username string, password string) (string, error) {
	return "user-1", nil
}

func (s *userServiceSuccessStub) CreateUser(req models.UserRequest) (*models.User, error) {
	user := models.User{
		Id:        "new-id",
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.users["new-id"] = user
	return &user, nil
}

func (s *userServiceSuccessStub) GetUsers() ([]models.UserResponse, error) {
	responses := make([]models.UserResponse, 0, len(s.users))
	for _, u := range s.users {
		responses = append(responses, u.ToResponse())
	}
	return responses, nil
}

func (s *userServiceSuccessStub) GetUserById(id string) (*models.UserResponse, error) {
	u, exists := s.users[id]
	if !exists {
		return nil, nil
	}
	resp := u.ToResponse()
	return &resp, nil
}

func (s *userServiceSuccessStub) UpdateUser(id string, req models.UserRequest) (*models.UserResponse, error) {
	u, exists := s.users[id]
	if !exists {
		return nil, nil
	}
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	s.users[id] = u
	resp := u.ToResponse()
	return &resp, nil
}

func (s *userServiceSuccessStub) DeleteUser(id string) error {
	if _, exists := s.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(s.users, id)
	return nil
}

type userServiceErrorStub struct{}

func (s *userServiceErrorStub) GetUserId(username string, password string) (string, error) {
	return "", nil
}
func (s *userServiceErrorStub) CreateUser(req models.UserRequest) (*models.User, error) {
	return nil, errors.New("username already exists")
}
func (s *userServiceErrorStub) GetUsers() ([]models.UserResponse, error) {
	return nil, errors.New("db error")
}
func (s *userServiceErrorStub) GetUserById(id string) (*models.UserResponse, error) {
	return nil, errors.New("db error")
}
func (s *userServiceErrorStub) UpdateUser(id string, req models.UserRequest) (*models.UserResponse, error) {
	return nil, errors.New("db error")
}
func (s *userServiceErrorStub) DeleteUser(id string) error {
	return errors.New("db error")
}

type userServiceEmptyStub struct{}

func (s *userServiceEmptyStub) GetUserId(username string, password string) (string, error) {
	return "", nil
}
func (s *userServiceEmptyStub) CreateUser(req models.UserRequest) (*models.User, error) {
	return nil, nil
}
func (s *userServiceEmptyStub) GetUsers() ([]models.UserResponse, error) {
	return []models.UserResponse{}, nil
}
func (s *userServiceEmptyStub) GetUserById(id string) (*models.UserResponse, error) {
	return nil, nil
}
func (s *userServiceEmptyStub) UpdateUser(id string, req models.UserRequest) (*models.UserResponse, error) {
	return nil, nil
}
func (s *userServiceEmptyStub) DeleteUser(id string) error {
	return errors.New("user not found")
}

func TestUserController_CreateUser_Success(t *testing.T) {
	svc := newUserServiceSuccessStub()
	ctrl := NewUserController(svc)

	e := echo.New()
	body, _ := json.Marshal(models.UserRequest{
		Username: "newuser", Password: "pass123",
		Email: "new@test.com", FirstName: "New", LastName: "User",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/user/", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ctrl.CreateUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp models.UserResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "newuser", resp.Username)
	}
}

func TestUserController_CreateUser_BadRequest(t *testing.T) {
	svc := newUserServiceSuccessStub()
	ctrl := NewUserController(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/user/", bytes.NewReader([]byte("{invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ctrl.CreateUser(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestUserController_CreateUser_ValidationError(t *testing.T) {
	svc := newUserServiceSuccessStub()
	ctrl := NewUserController(svc)

	e := echo.New()
	body, _ := json.Marshal(models.UserRequest{Password: "onlypassword"})
	req := httptest.NewRequest(http.MethodPost, "/api/user/", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ctrl.CreateUser(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestUserController_CreateUser_Conflict(t *testing.T) {
	svc := &userServiceErrorStub{}
	ctrl := NewUserController(svc)

	e := echo.New()
	body, _ := json.Marshal(models.UserRequest{
		Username: "existing", Password: "pass123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/user/", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ctrl.CreateUser(c)) {
		assert.Equal(t, http.StatusConflict, rec.Code)
	}
}

func TestUserController_GetUsers_Success(t *testing.T) {
	svc := newUserServiceSuccessStub()
	ctrl := NewUserController(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/user/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ctrl.GetUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var users []models.UserResponse
		err := json.Unmarshal(rec.Body.Bytes(), &users)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1)
	}
}

func TestUserController_GetUsers_Error(t *testing.T) {
	svc := &userServiceErrorStub{}
	ctrl := NewUserController(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/user/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ctrl.GetUsers(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}

func TestUserController_GetUserById_Found(t *testing.T) {
	svc := newUserServiceSuccessStub()
	ctrl := NewUserController(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/user/user-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues("user-1")

	if assert.NoError(t, ctrl.GetUserById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp models.UserResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "existing", resp.Username)
	}
}

func TestUserController_GetUserById_NotFound(t *testing.T) {
	svc := newUserServiceSuccessStub()
	ctrl := NewUserController(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/user/nonexistent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues("nonexistent")

	if assert.NoError(t, ctrl.GetUserById(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestUserController_UpdateUser_Success(t *testing.T) {
	svc := newUserServiceSuccessStub()
	ctrl := NewUserController(svc)

	e := echo.New()
	body, _ := json.Marshal(models.UserRequest{
		Username: "updated", Password: "pass123",
		FirstName: "Updated", LastName: "User",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/user/user-1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues("user-1")

	if assert.NoError(t, ctrl.UpdateUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestUserController_UpdateUser_NotFound(t *testing.T) {
	svc := newUserServiceSuccessStub()
	ctrl := NewUserController(svc)

	e := echo.New()
	body, _ := json.Marshal(models.UserRequest{
		Username: "any", Password: "pass123",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/user/nonexistent", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues("nonexistent")

	if assert.NoError(t, ctrl.UpdateUser(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestUserController_DeleteUser_Success(t *testing.T) {
	svc := newUserServiceSuccessStub()
	ctrl := NewUserController(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/user/user-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues("user-1")

	if assert.NoError(t, ctrl.DeleteUser(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}

func TestUserController_DeleteUser_NotFound(t *testing.T) {
	svc := newUserServiceSuccessStub()
	ctrl := NewUserController(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/user/nonexistent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues("nonexistent")

	if assert.NoError(t, ctrl.DeleteUser(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestUserController_DeleteUser_InternalError(t *testing.T) {
	svc := &userServiceErrorStub{}
	ctrl := NewUserController(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/user/any", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues("any")

	if assert.NoError(t, ctrl.DeleteUser(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}
