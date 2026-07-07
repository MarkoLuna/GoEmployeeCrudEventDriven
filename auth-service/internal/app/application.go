package app

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/MarkoLuna/AuthService/internal/config"
	"github.com/MarkoLuna/AuthService/internal/controllers"
	"github.com/MarkoLuna/AuthService/internal/repositories"
	"github.com/MarkoLuna/AuthService/internal/routes"
	"github.com/MarkoLuna/AuthService/internal/services"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
	"github.com/labstack/echo/v4"
)

type Application struct {
	EchoInstance   *echo.Echo
	DbConnection   *sql.DB
	UserRepository repositories.UserRepository
	UserService    services.UserService
	AuthStrategy   services.AuthStrategy
	OAuthController   controllers.OAuthController
	UserController     controllers.UserController
}

func (app *Application) LoadConfiguration() {
	config.EnableCORS(app.EchoInstance)
	config.ConfigureLogging()
	config.ConfigureLoggingMiddleware(app.EchoInstance)
}

func (app *Application) Address() string {
	port := utils.GetEnv("SERVER_PORT", "8082")
	host := utils.GetEnv("SERVER_HOST", "0.0.0.0")
	return host + ":" + port
}

func (app *Application) HandleRoutes() {
	routes.RegisterSwaggerRoute(app.EchoInstance)
	routes.RegisterHealthcheckRoute(app.EchoInstance)
	routes.RegisterAuthRoutes(app.EchoInstance, &app.OAuthController)
	routes.RegisterUserRoutes(app.EchoInstance, &app.UserController)
}

func (app *Application) StartServer() {
	app.HandleRoutes()
	address := app.Address()
	log.Println("Starting auth-service on:", address)
	log.Fatal(app.EchoInstance.Start(address))
}

func (app *Application) StartSecureServer() {
	app.HandleRoutes()
	address := app.Address()
	log.Println("Starting auth-service on:", address)

	path, _ := filepath.Abs("../../resources/ssl/cert.pem")
	certFile := utils.GetEnv("SERVER_SSL_CERT_FILE_PATH", path+"/resources/ssl/cert.pem")
	keyFile := utils.GetEnv("SERVER_SSL_KEY_FILE_PATH", path+"/resources/ssl/key.pem")

	if err := app.EchoInstance.StartTLS(address, certFile, keyFile); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (app *Application) Run() {
	serverSslEnabled := utils.GetEnv("SERVER_SSL_ENABLED", "false")
	sslEnabled, _ := strconv.ParseBool(serverSslEnabled)
	if sslEnabled {
		app.StartSecureServer()
	} else {
		app.StartServer()
	}
}
