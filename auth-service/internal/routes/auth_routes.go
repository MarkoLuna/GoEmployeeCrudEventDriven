package routes

import (
	"github.com/MarkoLuna/AuthService/internal/controllers"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RegisterHealthcheckRoute(echoInstance *echo.Echo) {
	echoInstance.GET("/healthcheck/", controllers.HealthCheckHandler)
}

func RegisterSwaggerRoute(echoInstance *echo.Echo) {
	echoInstance.GET("/swagger/*", echoSwagger.WrapHandler)
}

func RegisterAuthRoutes(echoInstance *echo.Echo, oauthController *controllers.OAuthController) {
	oauthGroup := echoInstance.Group("/oauth")
	oauthGroup.POST("/token", oauthController.TokenHandler)
	oauthGroup.GET("/userinfo", oauthController.GetUserInfo)
}

func RegisterUserRoutes(echoInstance *echo.Echo, userController *controllers.UserController) {
	userGroup := echoInstance.Group("/api/user")
	userGroup.POST("/", userController.CreateUser)
	userGroup.GET("/", userController.GetUsers)
	userGroup.GET("/:userId", userController.GetUserById)
	userGroup.PUT("/:userId", userController.UpdateUser)
	userGroup.DELETE("/:userId", userController.DeleteUser)
}
