package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/MarkoLuna/EmployeeConsumer/internal/config"
	"github.com/MarkoLuna/EmployeeConsumer/internal/controllers"
	"github.com/MarkoLuna/EmployeeConsumer/internal/repositories"
	"github.com/MarkoLuna/EmployeeConsumer/internal/routes"
	"github.com/MarkoLuna/EmployeeConsumer/internal/services"
	commonAuth "github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/services/auth"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
	"github.com/labstack/echo/v4"
)

type Application struct {
	EchoInstance                 *echo.Echo
	DbConnection                 *sql.DB
	EmployeeService              services.EmployeeService
	EmployeeKafkaConsumerService services.KafkaConsumerService
	EmployeeRepository           repositories.EmployeeRepository
	EmployeeController           controllers.EmployeeController
	ValidationClient             *commonAuth.TokenValidationClient
}

func (app *Application) LoadConfiguration() {
	config.EnableCORS(app.EchoInstance)
	config.ConfigureLogging()

	serverAuthEnabled := utils.GetEnv("OAUTH_ENABLED", "false")
	authEnabled, _ := strconv.ParseBool(serverAuthEnabled)
	config.NewAuthConfig(app.EchoInstance, authEnabled, config.DefaultSkippedPaths[:], app.ValidationClient)
}

func (app *Application) Address() string {
	port := utils.GetEnv("SERVER_PORT", "8081")
	host := utils.GetEnv("SERVER_HOST", "0.0.0.0")
	return host + ":" + port
}

func (app *Application) HandleRoutes() {
	routes.RegisterSwaggerRoute(app.EchoInstance)
	routes.RegisterHealthcheckRoute(app.EchoInstance)
	routes.RegisterEmployeeStoreRoutes(app.EchoInstance, &app.EmployeeController)
}

func (app *Application) StartServer(ctx context.Context) {
	app.HandleRoutes()
	address := app.Address()
	log.Println("Starting server on:", address)

	serverErr := make(chan error, 1)
	go func() {
		if err := app.EchoInstance.Start(address); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	case <-ctx.Done():
		log.Println("HTTP server: shutting down gracefully")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.EchoInstance.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server: shutdown error: %v", err)
		}
	}
}

func (app *Application) StartSecureServer(ctx context.Context) {
	app.HandleRoutes()
	address := app.Address()
	log.Println("Starting server on:", address)

	path, _ := filepath.Abs("../../resources/ssl/cert.pem")
	certFile := utils.GetEnv("SERVER_SSL_CERT_FILE_PATH", path+"/resources/ssl/cert.pem")
	keyFile := utils.GetEnv("SERVER_SSL_KEY_FILE_PATH", path+"/resources/ssl/key.pem")

	serverErr := make(chan error, 1)
	go func() {
		if err := app.EchoInstance.StartTLS(address, certFile, keyFile); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("HTTPS server error: %v", err)
		}
	case <-ctx.Done():
		log.Println("HTTPS server: shutting down gracefully")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.EchoInstance.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTPS server: shutdown error: %v", err)
		}
	}
}

func (app *Application) Run(ctx context.Context) {
	serverSslEnabled := utils.GetEnv("SERVER_SSL_ENABLED", "false")
	sslEnabled, _ := strconv.ParseBool(serverSslEnabled)
	if sslEnabled {
		app.StartSecureServer(ctx)
	} else {
		app.StartServer(ctx)
	}
}
