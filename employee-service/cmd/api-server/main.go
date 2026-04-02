package main

import (
	"github.com/MarkoLuna/EmployeeService/internal/app"
	"github.com/MarkoLuna/EmployeeService/internal/clients"
	"github.com/MarkoLuna/EmployeeService/internal/config"
	"github.com/MarkoLuna/EmployeeService/internal/controllers"
	"github.com/MarkoLuna/EmployeeService/internal/repositories"
	"github.com/MarkoLuna/EmployeeService/internal/services"
	"github.com/MarkoLuna/EmployeeService/internal/services/impl"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/services/auth"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
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
	defer App.DbConnection.Close()
	App.Run()
}

func ConfigureApp() {
	App.EchoInstance = echo.New()
	if App.DbConnection == nil {
		App.DbConnection = config.GetDB()
	}

	if App.EmployeeRepository == nil {
		App.EmployeeRepository = repositories.NewEmployeeRepository(App.DbConnection, true)
	}

	if App.EmployeeConsumerServiceClient == nil {
		httpClient := config.NewHttpClient()
		App.EmployeeConsumerServiceClient = clients.NewEmployeeConsumerServiceClient(*httpClient)
	}

	if App.KafkaProducerService == nil {
		kafkaConsumer, err := config.NewKafkaProducer()
		if err != nil {
			panic(err)
		}

		App.KafkaProducerService = impl.NewKafkaProducerService(kafkaConsumer)
	}

	App.EmployeeService = services.NewEmployeeService(App.EmployeeConsumerServiceClient, App.KafkaProducerService)
	App.EmployeeController = controllers.NewEmployeeController(App.EmployeeService)

	App.ClientService = services.NewClientService()
	App.UserService = services.NewUserService()

	oauthProvider := utils.GetEnv("OAUTH_PROVIDER", "local")
	if oauthProvider == "keycloak" {
		authServerURL := utils.GetEnv("KEYCLOAK_AUTH_SERVER_URL", "http://localhost:8082")
		realm := utils.GetEnv("KEYCLOAK_REALM", "dev")
		App.OAuthService = auth.NewKeycloakOAuthService(authServerURL, realm)
	} else {
		App.OAuthService = impl.NewLocalOAuthService()
	}

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
}
