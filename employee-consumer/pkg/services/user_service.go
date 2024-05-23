package services

import (
	"errors"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/utils"
)

var (
	userId            = utils.GetEnv("OAUTH_USER_ID", "000000")
	oauthUserName     = utils.GetEnv("OAUTH_USER_NAME", "user")
	oauthUserPassword = utils.GetEnv("OAUTH_USER_PASSWORD", "secret")
)

type UserService struct {
}

func NewUserService() UserService {
	return UserService{}
}

func (eSrv UserService) GetUserId(userName string, password string) (string, error) {

	if oauthUserName != userName {
		return "", errors.New("invalid username")
	}

	if oauthUserPassword != password {
		return "", errors.New("invalid password")
	}

	return userId, nil
}
