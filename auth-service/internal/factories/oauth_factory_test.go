package factories

import (
	"testing"

	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/services/auth"
	"github.com/stretchr/testify/assert"
)

func TestGetOAuthProviderFactory_Keycloak(t *testing.T) {
	factory := GetOAuthProviderFactory("keycloak")
	assert.NotNil(t, factory)
	svc := factory()
	assert.NotNil(t, svc)
}

func TestGetOAuthProviderFactory_Local(t *testing.T) {
	factory := GetOAuthProviderFactory("local")
	assert.NotNil(t, factory)
	svc := factory()
	assert.NotNil(t, svc)

	_, ok := svc.(auth.OAuthService)
	assert.True(t, ok, "local provider should implement OAuthService")
}

func TestGetOAuthProviderFactory_UnknownProvider(t *testing.T) {
	factory := GetOAuthProviderFactory("nonexistent")
	assert.NotNil(t, factory)
	svc := factory()
	assert.NotNil(t, svc)

	_, ok := svc.(auth.OAuthService)
	assert.True(t, ok, "unknown provider should fallback to local")
}
