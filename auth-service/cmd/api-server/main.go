package main

import (
	"log"

	"github.com/MarkoLuna/AuthService/internal/app"
	"github.com/MarkoLuna/AuthService/internal/config"
	"github.com/MarkoLuna/AuthService/internal/controllers"
	"github.com/MarkoLuna/AuthService/internal/repositories"
	"github.com/MarkoLuna/AuthService/internal/services/impl"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
	"github.com/labstack/echo/v4"

	_ "github.com/MarkoLuna/AuthService/docs"
)

var (
	App = app.Application{}
)

// @title Auth Service API
// @version 1.0
// @description Authentication and User Management Service

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email josemarcosluna9@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
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
	config.InitDatabase()
	if App.DbConnection == nil {
		App.DbConnection = config.GetDB()
	}

	if App.UserRepository == nil {
		App.UserRepository = repositories.NewUserRepository(App.DbConnection, true)
	}

	App.UserService = impl.NewUserService(App.UserRepository)

	oauthProvider := utils.GetEnv("OAUTH_PROVIDER", "local")
	log.Println("OAuth Provider: ", oauthProvider)

	switch oauthProvider {
	case "keycloak":
		App.AuthStrategy = impl.NewKeycloakAuthStrategy(
			utils.GetEnv("KEYCLOAK_AUTH_SERVER_URL", "http://localhost:8082"),
			utils.GetEnv("KEYCLOAK_REALM", "dev"),
		)
	default:
		App.AuthStrategy = impl.NewLocalAuthStrategy(
			impl.NewClientService(App.UserRepository),
			App.UserService,
			impl.NewLocalOAuthService(),
		)
	}

	App.OAuthController = controllers.NewOAuthController(App.AuthStrategy)
	App.UserController = controllers.NewUserController(App.UserService)

	App.LoadConfiguration()
}
