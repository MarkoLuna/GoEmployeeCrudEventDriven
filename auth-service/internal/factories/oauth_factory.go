package factories

import (
	"log"

	"github.com/MarkoLuna/AuthService/internal/services/impl"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/services/auth"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
)

var OauthProviderFactories = map[string]func() auth.OAuthService{
	"keycloak": func() auth.OAuthService {
		log.Println("Using Keycloak for OAuth")
		authServerURL := utils.GetEnv("KEYCLOAK_AUTH_SERVER_URL", "http://localhost:8082")
		realm := utils.GetEnv("KEYCLOAK_REALM", "dev")
		return auth.NewKeycloakOAuthService(authServerURL, realm)
	},
	"local": func() auth.OAuthService {
		log.Println("Using Local OAuth")
		return impl.NewLocalOAuthService()
	},
}

func GetOAuthProviderFactory(provider string) func() auth.OAuthService {
	factory, exists := OauthProviderFactories[provider]
	if !exists {
		log.Printf("Warning: OAuth Provider '%s' not found, falling back to local\n", provider)
		factory = OauthProviderFactories["local"]
	}
	return factory
}
