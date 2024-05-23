package main

import (
	"github.com/MarkoLuna/EmployeeConsumer/pkg/app"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/config"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/controllers"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/repositories"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/services"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/services/impl"
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

	App.EmployeeService = services.NewEmployeeService(App.EmployeeRepository)
	App.EmployeeController = controllers.NewEmployeeController(App.EmployeeService)

	App.ClientService = services.NewClientService()
	App.UserService = services.NewUserService()

	if App.OAuthService == nil {
		App.OAuthService = impl.NewOAuthService()
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
