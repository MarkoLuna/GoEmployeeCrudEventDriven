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
	commonAuth "github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/services/auth"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
	"github.com/labstack/echo/v4"
)

type Application struct {
	EchoInstance                  *echo.Echo
	EmployeeService               services.EmployeeService
	EmployeeController                   controllers.EmployeeController
	EmployeeConsumerServiceClientBuilder *clients.EmployeeConsumerServiceClientBuilder
	KafkaProducerService                 services.KafkaProducerService
	ValidationClient                     *commonAuth.TokenValidationClient
}

func (app *Application) LoadConfiguration() {
	config.EnableCORS(app.EchoInstance)
	config.ConfigureLogging()

	serverAuthEnabled := utils.GetEnv("OAUTH_ENABLED", "false")
	authEnabled, _ := strconv.ParseBool(serverAuthEnabled)
	config.NewAuthConfig(app.EchoInstance, authEnabled, config.DefaultSkippedPaths[:], app.ValidationClient)
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
