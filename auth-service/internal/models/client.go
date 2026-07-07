package models

import "time"

type Client struct {
	Id           string    `json:"id"`
	ClientId     string    `json:"clientId"`
	ClientSecret string    `json:"-"`
	Enabled      bool      `json:"enabled"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type ClientRequest struct {
	ClientId     string `json:"clientId" validate:"required"`
	ClientSecret string `json:"clientSecret" validate:"required"`
}

type ClientResponse struct {
	Id        string    `json:"id"`
	ClientId  string    `json:"clientId"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c Client) ToResponse() ClientResponse {
	return ClientResponse{
		Id:        c.Id,
		ClientId:  c.ClientId,
		Enabled:   c.Enabled,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
