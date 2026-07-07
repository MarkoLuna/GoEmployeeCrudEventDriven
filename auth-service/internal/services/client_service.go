package services

type ClientService interface {
	IsValidClientCredentials(clientId string, clientSecret string) (bool, error)
}
