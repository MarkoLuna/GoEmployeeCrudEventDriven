package main

import (
	"log"

	"github.com/MarkoLuna/EmployeeService/internal/app"
	"github.com/MarkoLuna/EmployeeService/internal/clients"
	"github.com/MarkoLuna/EmployeeService/internal/config"
	"github.com/MarkoLuna/EmployeeService/internal/controllers"
	"github.com/MarkoLuna/EmployeeService/internal/services"
	"github.com/MarkoLuna/EmployeeService/internal/services/impl"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/services/auth"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
	"github.com/labstack/echo/v4"

	_ "github.com/MarkoLuna/EmployeeService/docs"
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
	ConfigureApp()
	App.Run()
}

func ConfigureApp() {
	App.EchoInstance = echo.New()
	if App.EmployeeConsumerServiceClientBuilder == nil {
		httpClient := config.NewHttpClient()
		App.EmployeeConsumerServiceClientBuilder = clients.NewEmployeeConsumerServiceClientBuilder().WithHttpClient(*httpClient)
	}

	if App.KafkaProducerService == nil {
		kafkaConsumer, err := config.NewKafkaProducer()
		if err != nil {
			panic(err)
		}

		App.KafkaProducerService = impl.NewKafkaProducerService(kafkaConsumer)
	}

	App.EmployeeService = services.NewEmployeeService(App.EmployeeConsumerServiceClientBuilder, App.KafkaProducerService)
	App.EmployeeController = controllers.NewEmployeeController(App.EmployeeService)

	App.ValidationClient = auth.NewTokenValidationClient(
		utils.GetEnv("AUTH_SERVICE_URL", "http://localhost:8082"),
	)

	log.Println("OAuth handled by external auth-service")

	App.LoadConfiguration()
}
