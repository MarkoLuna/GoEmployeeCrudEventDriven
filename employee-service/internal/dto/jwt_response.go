package dto

import (
	"encoding/json"
	"log"
)

type JWTResponse struct {
	AccessToken  string `json:"access_token" swaggertype:"string" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjbGllbnQiLCJleHAiOjE2NjIxNTY3MDUsInN1YiI6IjAwMDAwMCJ9.xqHSNgbzZWFCmkMw48syhVJvQkyvnnM7__Rk915EMv2Di2kdIFiZJwWIt9RciD2jKgyBB-Usei3wEwzxuHhLgQ"`
	RefreshToken string `json:"refresh_token" swaggertype:"string" example:"NMU1NDCWODYTNGZIMY01YMVLLTLLMGETMJU2ZDNLNTJIMGI5"`
	ExpiresIn    int64  `json:"expires_in" swaggertype:"integer" example:"120"`
	Scope        string `json:"scope" swaggertype:"string" example:"all"`
	TokenType    string `json:"token_type" swaggertype:"string" example:"Bearer"`
}

func (e JWTResponse) ToString() string {
	out, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}
