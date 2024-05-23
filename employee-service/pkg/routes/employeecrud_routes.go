package routes

import (
	"github.com/MarkoLuna/EmployeeService/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterEmployeeStoreRoutes(echoInstance *echo.Echo, employeeController *controllers.EmployeeController) {

	employeeGroup := echoInstance.Group("/api/employee")
	employeeGroup.POST("/", employeeController.CreateEmployee)
	employeeGroup.GET("/", employeeController.GetEmployees)
	employeeGroup.GET("/:employeeId", employeeController.GetEmployeeById)
	employeeGroup.PUT("/:employeeId", employeeController.UpdateEmployee)
	employeeGroup.DELETE("/:employeeId", employeeController.DeleteEmployee)
}

func RegisterHealthcheckRoute(echoInstance *echo.Echo) {
	echoInstance.GET("/healthcheck/", controllers.HealthCheckHandler)
}

func RegisterOAuthRoutes(echoInstance *echo.Echo, oauthController *controllers.OAuthController) {
	oauthGroup := echoInstance.Group("/oauth")
	oauthGroup.POST("/token", oauthController.TokenHandler)
	oauthGroup.GET("/userinfo", oauthController.GetUserInfo)
}
