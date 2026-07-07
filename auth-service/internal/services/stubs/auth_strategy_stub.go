package stubs

import (
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/dto"
	"github.com/labstack/echo/v4"
)

type AuthStrategyStub struct {
}

func NewAuthStrategyStub() AuthStrategyStub {
	return AuthStrategyStub{}
}

func (s AuthStrategyStub) HandleTokenGeneration(c echo.Context) (*dto.JWTResponse, error) {
	access := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjNmNlY2U1MyIsImV4cCI6MTY0Mjc5MTUzNiwic3ViIjoiMDAwMDAwIn0.SA49Q2UZvzf7dgZmvzNTaBjF1aYP821iXZje2pxK1KgvjZlQNOmQQ1B1duxfDkXeIWUfbFi2dkzlXx4GcWOVeg"
	refresh := "ZJVLYTRINZUTZJNLMY01MZLLLWFJMMMTYMM3Y2YZNTQ1MWM2"

	return &dto.JWTResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    120,
		Scope:        "all",
		TokenType:    "Bearer",
	}, nil
}

func (s AuthStrategyStub) GetTokenClaims(accessToken string) (map[string]string, error) {
	return make(map[string]string), nil
}
