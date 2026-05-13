package app

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/MarkoLuna/EmployeeService/internal/clients"
	"github.com/MarkoLuna/EmployeeService/internal/config"
	"github.com/MarkoLuna/EmployeeService/internal/controllers"
	"github.com/MarkoLuna/EmployeeService/internal/routes"
	"github.com/MarkoLuna/EmployeeService/internal/services"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
	"github.com/labstack/echo/v4"
)

type Application struct {
	EchoInstance                  *echo.Echo
	OAuthService                  services.OAuthService
	EmployeeService               services.EmployeeService
	ClientService                 services.ClientService
	UserService                   services.UserService
	EmployeeController                   controllers.EmployeeController
	OAuthController                      controllers.OAuthController
	EmployeeConsumerServiceClientBuilder *clients.EmployeeConsumerServiceClientBuilder
	KafkaProducerService                 services.KafkaProducerService
}

func (app *Application) LoadConfiguration() {
	config.EnableCORS(app.EchoInstance)
	config.ConfigureLogging()

	server_auth_enabled := utils.GetEnv("OAUTH_ENABLED", "false")
	auth_enabled, _ := strconv.ParseBool(server_auth_enabled)
	config.NewAuthConfig(app.EchoInstance, auth_enabled, config.DefaultSkippedPaths[:], app.OAuthService)
}

func (app *Application) Address() string {
	port := utils.GetEnv("SERVER_PORT", "8080")
	host := utils.GetEnv("SERVER_HOST", "0.0.0.0")

	return host + ":" + port
}

func (app *Application) HandleRoutes() {
	routes.RegisterSwaggerRoute(app.EchoInstance)
	routes.RegisterHealthcheckRoute(app.EchoInstance)
	routes.RegisterEmployeeStoreRoutes(app.EchoInstance, &app.EmployeeController)
	routes.RegisterOAuthRoutes(app.EchoInstance, &app.OAuthController)
}

func (app *Application) StartServer() {
	app.HandleRoutes()
	address := app.Address()
	log.Println("Starting server on:", address)
	log.Fatal(app.EchoInstance.Start(address))
}

func (app *Application) StartSecureServer() {
	app.HandleRoutes()
	address := app.Address()
	log.Println("Starting server on:", address)

	// path := "/Users/marcos.luna/go-projects/GoEmployeeCrudEventDriven/EmployeeService"
	path, _ := filepath.Abs("../../resources/ssl/cert.pem")
	certFile := utils.GetEnv("SERVER_SSL_CERT_FILE_PATH", path+"/resources/ssl/cert.pem")
	keyFile := utils.GetEnv("SERVER_SSL_KEY_FILE_PATH", path+"/resources/ssl/key.pem")

	if err := app.EchoInstance.StartTLS(address, certFile, keyFile); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (app *Application) Run() {
	server_ssl_enabled := utils.GetEnv("SERVER_SSL_ENABLED", "false")
	ssl_enabled, _ := strconv.ParseBool(server_ssl_enabled)
	if ssl_enabled {
		app.StartSecureServer()
	} else {
		app.StartServer()
	}
}
