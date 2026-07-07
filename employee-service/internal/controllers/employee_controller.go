package controllers

import (
	"log"
	"net/http"

	"github.com/MarkoLuna/EmployeeService/internal/models"
	"github.com/MarkoLuna/EmployeeService/internal/services"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

var NewEmployee models.Employee

type EmployeeController struct {
	employeeService services.EmployeeService
}

func NewEmployeeController(employeeService services.EmployeeService) EmployeeController {
	return EmployeeController{employeeService}
}

func (eCtrl EmployeeController) getJwtToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if len(authHeader) >= 7 && authHeader[0:7] == "Bearer " {
		return authHeader[7:]
	}
	return authHeader
}

// CreateEmployee EmployeeApi
// @Tags 	EmployeeApi
// @Summary create-employee
// @Description Add a new employee to the database
// @Accept  json
// @Produce  json
// @Param   employee-details      body dto.EmployeeRequest true  "Some ID"
// @Success 200 {string} string	"ok"
// @Failure 400 {object} error "Invalid request!!"
// @Security ApiKeyAuth
// @Router /api/employee/ [post]
func (eCtrl EmployeeController) CreateEmployee(c echo.Context) error {
	employee := dto.EmployeeRequest{}
	if err := c.Bind(&employee); err != nil {
		return err
	}

	v := utils.CreateValidator()
	err := v.Struct(employee)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Println(e)
		}
		return c.String(http.StatusBadRequest, "")
	}

	jwt := eCtrl.getJwtToken(c)
	_, err = eCtrl.employeeService.CreateEmployee(c.Request().Context(), jwt, employee)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, map[string]string{
		"status":  "accepted",
		"message": "Employee creation request submitted",
	})
}

// GetEmployees EmployeeApi
// @Tags 	EmployeeApi
// @Summary get-employees
// @Description Get employees from the database
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.Employee	"ok"
// @Failure 400 {object} error "Invalid request!!"
// @Security ApiKeyAuth
// @Router /api/employee/ [get]
func (eCtrl EmployeeController) GetEmployees(c echo.Context) error {
	jwt := eCtrl.getJwtToken(c)
	newEmployees, err := eCtrl.employeeService.GetEmployees(c.Request().Context(), jwt)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, newEmployees)
}

// GetEmployee EmployeeApi
// @Tags 	EmployeeApi
// @Summary get-employee
// @Description Get employee from the database
// @Accept  json
// @Produce  json
// @Param   employeeId      path string true  "Employee ID"
// @Success 200 {object} models.Employee	"ok"
// @Failure 400 {object} error "Invalid request!!"
// @Security ApiKeyAuth
// @Router /api/employee/{employeeId} [get]
func (eCtrl EmployeeController) GetEmployeeById(c echo.Context) error {
	jwt := eCtrl.getJwtToken(c)
	employeeId := c.Param("employeeId")
	EmployeeDetails, err := eCtrl.employeeService.GetEmployeeById(c.Request().Context(), jwt, employeeId)
	if err == nil {
		return c.JSON(http.StatusOK, EmployeeDetails)
	} else {
		return c.String(http.StatusNotFound, err.Error())
	}
}

// UpdateEmployee EmployeeApi
// @Tags 	EmployeeApi
// @Summary update-employee
// @Description Update employee
// @Accept  json
// @Produce  json
// @Param   employeeId      path string true  "Employee ID"
// @Param   employee-details      body dto.EmployeeRequest true  "Some ID"
// @Success 200 {object} models.Employee	"ok"
// @Failure 400 {object} error "Invalid request!!"
// @Security ApiKeyAuth
// @Router /api/employee/{employeeId} [put]
func (eCtrl EmployeeController) UpdateEmployee(c echo.Context) error {
	var updateEmployee = dto.EmployeeRequest{}
	if err := c.Bind(&updateEmployee); err != nil {
		return err
	}

	log.Println("employee: " + updateEmployee.ToString())

	v := utils.CreateValidator()
	err := v.Struct(updateEmployee)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Println(e)
		}
		return c.String(http.StatusNotFound, err.Error())
	}

	jwt := eCtrl.getJwtToken(c)
	employeeId := c.Param("employeeId")
	_, err = eCtrl.employeeService.UpdateEmployee(c.Request().Context(), jwt, employeeId, updateEmployee)
	if err == nil {
		return c.JSON(http.StatusAccepted, map[string]string{
			"status":  "accepted",
			"message": "Employee update request submitted",
		})
	} else {
		return c.String(http.StatusNotFound, err.Error())
	}
}

// DeleteEmployee EmployeeApi
// @Tags 	EmployeeApi
// @Summary delete-employee
// @Description Delete employee from the database
// @Accept  json
// @Produce  json
// @Param   employeeId      path string true  "Employee ID"
// @Success 200 {string} string	"ok"
// @Failure 400 {object} error "Invalid request!!"
// @Security ApiKeyAuth
// @Router /api/employee/{employeeId} [delete]
func (eCtrl EmployeeController) DeleteEmployee(c echo.Context) error {
	jwt := eCtrl.getJwtToken(c)
	employeeId := c.Param("employeeId")

	err := eCtrl.employeeService.DeleteEmployeeById(c.Request().Context(), jwt, employeeId)
	if err == nil {
		return c.JSON(http.StatusAccepted, map[string]string{
			"status":  "accepted",
			"message": "Employee deletion request submitted",
		})
	} else {
		return c.String(http.StatusInternalServerError, err.Error())
	}

}
