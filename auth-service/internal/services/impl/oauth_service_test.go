package impl

import (
	"testing"

	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/stretchr/testify/assert"
)

func TestLocalOAuthServiceImpl_HandleTokenGeneration_Success(t *testing.T) {
	svc := NewLocalOAuthService()

	resp, err := svc.HandleTokenGeneration("test-client", "test-secret", "user-1")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, int64(120), resp.ExpiresIn)
	assert.Equal(t, "all", resp.Scope)
	assert.Equal(t, "Bearer", resp.TokenType)
}

func TestLocalOAuthServiceImpl_ParseToken_Valid(t *testing.T) {
	svc := NewLocalOAuthService()

	resp, err := svc.HandleTokenGeneration("client", "secret", "user-1")
	assert.NoError(t, err)

	token, err := svc.ParseToken(resp.AccessToken)
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.True(t, token.Valid)
}

func TestLocalOAuthServiceImpl_ParseToken_Invalid(t *testing.T) {
	svc := NewLocalOAuthService()

	token, err := svc.ParseToken("invalid-token-string")
	assert.Error(t, err)
	assert.Nil(t, token)
}

func TestLocalOAuthServiceImpl_GetTokenClaims_Success(t *testing.T) {
	svc := NewLocalOAuthService()

	resp, err := svc.HandleTokenGeneration("my-client", "my-secret", "user-123")
	assert.NoError(t, err)

	claims, err := svc.GetTokenClaims(resp.AccessToken)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user-123", claims["subject"])
}

func TestLocalOAuthServiceImpl_GetTokenClaims_Invalid(t *testing.T) {
	svc := NewLocalOAuthService()

	claims, err := svc.GetTokenClaims("garbage-token")
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestLocalOAuthServiceImpl_IsValidToken_Valid(t *testing.T) {
	svc := NewLocalOAuthService()

	resp, err := svc.HandleTokenGeneration("c", "s", "u")
	assert.NoError(t, err)

	valid, err := svc.IsValidToken(resp.AccessToken)
	assert.NoError(t, err)
	assert.True(t, valid)
}

func TestLocalOAuthServiceImpl_IsValidToken_Invalid(t *testing.T) {
	svc := NewLocalOAuthService()

	valid, err := svc.IsValidToken("bad-token")
	assert.Error(t, err)
	assert.False(t, valid)
}

func TestLocalOAuthServiceImpl_HandleTokenGeneration_ReturnsValidJWTResponse(t *testing.T) {
	svc := NewLocalOAuthService()

	resp, err := svc.HandleTokenGeneration("c", "s", "u")
	assert.NoError(t, err)

	var expected dto.JWTResponse
	expected = resp
	assert.Equal(t, expected.AccessToken, resp.AccessToken)
	assert.Equal(t, expected.RefreshToken, resp.RefreshToken)
	assert.Equal(t, expected.TokenType, resp.TokenType)
}
