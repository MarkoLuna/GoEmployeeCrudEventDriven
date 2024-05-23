package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/dto"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/utils"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var (
	signingKey = utils.GetEnv("OAUTH_SIGNING_KEY", "00000000")
)

type OAuthService struct {
}

func NewOAuthService() OAuthService {
	return OAuthService{}
}

func (eSrv OAuthService) HandleTokenGeneration(clientId string, clientSecret string, userId string) (dto.JWTResponse, error) {

	data := &oauth2.GenerateBasic{
		Client: &models.Client{
			ID:     clientId,
			Secret: clientSecret,
		},
		UserID: userId,
		TokenInfo: &models.Token{
			AccessCreateAt:  time.Now(),
			AccessExpiresIn: time.Second * 120,
		},
	}

	gen := generates.NewJWTAccessGenerate("", []byte(signingKey), jwt.SigningMethodHS512)
	access, refresh, err := gen.Token(context.Background(), data, true)

	if err != nil {
		return dto.JWTResponse{}, err
	}

	jWTResponse := dto.JWTResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int64(120), // 2 min
		Scope:        "all",
		TokenType:    "Bearer",
	}

	return jWTResponse, nil
}

func (oauthService OAuthService) ParseToken(accessToken string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(accessToken, &generates.JWTAccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("parse error")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (oauthService OAuthService) GetTokenClaims(accessToken string) (map[string]string, error) {

	token, err := oauthService.ParseToken(accessToken)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*generates.JWTAccessClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	dataMap := make(map[string]string)
	dataMap["subject"] = claims.Subject
	dataMap["audience"] = claims.Audience
	dataMap["id"] = claims.Id
	dataMap["issuer"] = claims.Issuer

	return dataMap, nil
}

func (oauthService OAuthService) IsAuthenticated(c echo.Context) (bool, error) {
	accessToken, _ := utils.GetBearerAuth(c.Request().Header)
	ok, err := oauthService.IsValidToken(accessToken)
	return ok && err != nil, err
}

func (oauthService OAuthService) IsValidToken(accessToken string) (bool, error) {

	// Parse and verify jwt access token
	token, err := jwt.ParseWithClaims(accessToken, &generates.JWTAccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt parse error")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return false, err
	}

	_, ok2 := token.Claims.(*generates.JWTAccessClaims)
	if !ok2 || !token.Valid {
		return false, err
	}
	return true, nil
}
