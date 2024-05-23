package controllers

import (
	"log"
	"net/http"

	"github.com/MarkoLuna/EmployeeService/pkg/dto"
	"github.com/MarkoLuna/EmployeeService/pkg/models"
	"github.com/MarkoLuna/EmployeeService/pkg/services"
	"github.com/MarkoLuna/EmployeeService/pkg/utils"
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

	e, err := eCtrl.employeeService.CreateEmployee(employee)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, e)
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
	newEmployees, err := eCtrl.employeeService.GetEmployees()
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
	employeeId := c.Param("employeeId")
	EmployeeDetails, err := eCtrl.employeeService.GetEmployeeById(employeeId)
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

	employeeId := c.Param("employeeId")
	employeeDetails, err := eCtrl.employeeService.UpdateEmployee(employeeId, updateEmployee)
	if err == nil {
		return c.JSON(http.StatusOK, employeeDetails)
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
	employeeId := c.Param("employeeId")

	err := eCtrl.employeeService.DeleteEmployeeById(employeeId)
	if err == nil {
		return c.String(http.StatusOK, "")
	} else {
		return c.String(http.StatusNotFound, err.Error())
	}

}
