package stubs

import (
	"github.com/MarkoLuna/EmployeeConsumer/pkg/dto"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type OAuthServiceStub struct {
}

func NewOAuthServiceStub() OAuthServiceStub {
	return OAuthServiceStub{}
}

func (eSrv OAuthServiceStub) HandleTokenGeneration(clientId string, clientSecret string, userId string) (dto.JWTResponse, error) {

	access := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9" + "." + "eyJhdWQiOiJjNmNlY2U1MyIsImV4cCI6MTY0Mjc5MTUzNiwic3ViIjoiMDAwMDAwIn0" + "." + "SA49Q2UZvzf7dgZmvzNTaBjF1aYP821iXZje2pxK1KgvjZlQNOmQQ1B1duxfDkXeIWUfbFi2dkzlXx4GcWOVeg"
	refresh := "ZJVLYTRINZUTZJNLMY01MZLLLWFJMMMTYMM3Y2YZNTQ1MWM2"

	jWTResponse := dto.JWTResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int64(120), // 2 min
		Scope:        "all",
		TokenType:    "Bearer",
	}

	return jWTResponse, nil
}

func (oauthService OAuthServiceStub) ParseToken(accessToken string) (*jwt.Token, error) {

	token := jwt.Token{}
	return &token, nil
}

func (oauthService OAuthServiceStub) IsAuthenticated(c echo.Context) (bool, error) {
	return true, nil
}

func (oauthService OAuthServiceStub) IsValidToken(accessToken string) (bool, error) {
	return true, nil
}

func (oauthService OAuthServiceStub) GetTokenClaims(accessToken string) (map[string]string, error) {
	claims := make(map[string]string)
	return claims, nil
}
