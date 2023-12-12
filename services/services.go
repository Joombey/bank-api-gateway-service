package services

import "farukh.go/api-gateway/models"

type ProfileService interface {
	TransferRequest(from int, to int, value float32) ([]models.ValueResponse, error)
	GetProfileRequest(id int) (models.ProfileResponse, error)
	CreateProfileRequest(username, password string) (models.ProfileResponse, error)
}

type KeycloakService interface {
	CheckToken() (bool, error)
	CheckRole(string) (bool, error)
	Login(username, password string) (string, error)
	RefreshToken() (models.Token, error)
	DeteUser(username string)
	RegisterUser(username, password string) error
}
