package handlers

import (
	"farukh.go/api-gateway/models"
	"farukh.go/api-gateway/di"
)

var repo = di.GetContiner().Profile
var cloak = di.GetContiner().Keycloak

func GetProfileHandler(id int) (models.ProfileResponse, error) {
	return repo.GetProfileRequest(id)
}

func Transfer(token models.Token, from int, to int, value float32) ([]models.ValueResponse, bool, error) {
	isValid, err := cloak.CheckToken()
	if err != nil {
		return nil, isValid, err
	} else if !isValid {
		return nil, true, nil
	}
	
	isValid, err = cloak.CheckRole("CARD_OWNER")
	if err != nil {
		return nil, isValid, err
	} else if !isValid {
		return nil, true, nil
	}
	

	response, err := repo.TransferRequest(from, to, value)
	if err != nil {
		return nil, false, err
	}

	return response, false, nil
}

func CreateUser(username, password string) {
	cloak.RegisterUser(username, password)
}