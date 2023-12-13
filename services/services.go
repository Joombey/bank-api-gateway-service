package services

import "farukh.go/api-gateway/models"

type ProfileService interface {
	TransferRequest(from int, to int, value float32) ([]models.ValueResponse, *models.ErrorBody)
	GetProfileRequest(id int) (models.ProfileResponse, *models.ErrorBody)
	CreateProfileRequest(username string) (models.ProfileResponse, *models.ErrorBody)
	Delete(id int) (models.ProfileResponse, *models.ErrorBody)
}
