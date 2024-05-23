package services

import (
	"errors"

	"github.com/MarkoLuna/EmployeeService/pkg/utils"
)

var (
	clientId     = utils.GetEnv("OAUTH_CLIENT_ID", "client")
	clientSecret = utils.GetEnv("OAUTH_CLIENT_SECRET", "password")
)

type ClientService struct {
}

func NewClientService() ClientService {
	return ClientService{}
}

func (eSrv ClientService) IsValidClientCredentials(client string, password string) (bool, error) {

	if clientId != client {
		return false, errors.New("invalid client id")
	}

	if clientSecret != password {
		return false, errors.New("invalid client secret")
	}

	return true, nil
}
