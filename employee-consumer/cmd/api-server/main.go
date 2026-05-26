package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MarkoLuna/EmployeeConsumer/internal/app"
	"github.com/MarkoLuna/EmployeeConsumer/internal/config"
	"github.com/MarkoLuna/EmployeeConsumer/internal/controllers"
	"github.com/MarkoLuna/EmployeeConsumer/internal/factories"
	"github.com/MarkoLuna/EmployeeConsumer/internal/repositories"
	"github.com/MarkoLuna/EmployeeConsumer/internal/services"
	"github.com/MarkoLuna/EmployeeConsumer/internal/services/impl"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/labstack/echo/v4"

	_ "github.com/MarkoLuna/EmployeeConsumer/docs"
)

var (
	App = app.Application{}
)

// @title Employee Crud API
// @version 1.0
// @description This app is responsable for a CRUD for Employees.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email josemarcosluna9@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Root context cancelled on SIGINT / SIGTERM so all components shut down together.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Shutdown coordinator manages graceful shutdown of background services.
	// This follows Dependency Inversion Principle by depending on abstraction.
	shutdownCoordinator := app.NewWaitGroupCoordinator()

	ConfigureApp(ctx, shutdownCoordinator)
	defer App.DbConnection.Close()

	// Block until the HTTP server exits (which itself waits for the signal).
	App.Run(ctx)

	// Wait for all background services to complete graceful shutdown.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := shutdownCoordinator.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown coordinator error: %v", err)
	}
	log.Println("Application shutdown complete")
}

func ConfigureApp(ctx context.Context, coordinator app.ShutdownCoordinator) {
	App.EchoInstance = echo.New()
	if App.DbConnection == nil {
		App.DbConnection = config.GetDB()
	}

	if App.EmployeeRepository == nil {
		App.EmployeeRepository = repositories.NewEmployeeRepository(App.DbConnection, true)
	}

	App.EmployeeService = impl.NewEmployeeService(App.EmployeeRepository)
	App.EmployeeController = controllers.NewEmployeeController(App.EmployeeService)
	kafkaConsumer, err := config.NewKafkaConsumer()
	if err != nil {
		log.Fatal("Unable to initialize kafka consumer due to Error", err)
	}

	kafkaProducer, err := config.NewKafkaProducer()
	if err != nil {
		log.Fatal("Unable to initialize kafka producer due to Error", err)
	}

	App.EmployeeKafkaConsumerService = impl.NewKafkaConsumerService(kafkaConsumer, kafkaProducer, App.EmployeeService)

	App.ClientService = services.NewClientService()
	App.UserService = services.NewUserService()

	oauthProvider := utils.GetEnv("OAUTH_PROVIDER", "local")
	log.Println("OAuth Provider: ", oauthProvider)

	authProviderFactory := factories.GetOAuthProviderFactory(oauthProvider)
	App.OAuthService = authProviderFactory()

	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	oauthServer := server.NewDefaultServer(manager)
	oauthServer.SetAllowGetAccessRequest(true)
	oauthServer.SetClientInfoHandler(server.ClientFormHandler)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	App.OAuthController = controllers.NewOAuthController(oauthServer, App.OAuthService, App.ClientService, App.UserService)

	App.LoadConfiguration()

	// Launch the Kafka consumer in its own goroutine using the coordinator.
	// This follows Single Responsibility Principle by separating lifecycle management.
	coordinator.Start(ctx, "kafka-consumer", App.EmployeeKafkaConsumerService.Listen)
}
