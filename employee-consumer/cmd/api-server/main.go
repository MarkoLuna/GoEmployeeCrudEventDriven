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
	"github.com/MarkoLuna/EmployeeConsumer/internal/repositories"
	"github.com/MarkoLuna/EmployeeConsumer/internal/services/impl"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/services/auth"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownCoordinator := app.NewWaitGroupCoordinator()

	ConfigureApp(ctx, shutdownCoordinator)
	defer App.DbConnection.Close()

	App.Run(ctx)

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

	App.ValidationClient = auth.NewTokenValidationClient(
		utils.GetEnv("AUTH_SERVICE_URL", "http://localhost:8082"),
	)

	log.Println("OAuth handled by external auth-service")

	App.LoadConfiguration()

	coordinator.Start(ctx, "kafka-consumer", App.EmployeeKafkaConsumerService.Listen)
}
