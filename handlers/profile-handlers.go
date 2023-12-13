package handlers

import (
	"net/http"

	cts "farukh.go/api-gateway/constants"
	"farukh.go/api-gateway/di"
	kk "farukh.go/api-gateway/keycloak"
	"farukh.go/api-gateway/models"
)

var repo = di.GetContiner().Profile

func GetProfileHandler(id int) (any, *models.ErrorBody) {
	return repo.GetProfileRequest(id)
}

func Login(username, password string) (models.Token, error) {
	return kk.LoginUser(username, password)
}

func Transfer(username string, token models.Token, from int, to int, value float32) (any, *models.ErrorBody) {
	newToken, errBody := preprocessUser(username, cts.RoleCardOwner, token)
	if errBody != nil {
		return newToken, errBody
	}

	serviceResponse, errBody := repo.TransferRequest(from, to, value)
	response := models.ReponseFrameWithError{Body: serviceResponse, NewToken: *newToken, Err: *errBody}

	return response, nil
}

func BlockUser(id int, username, target string, token models.Token) (any, *models.ErrorBody) {
	newToken, errBody := preprocessUser(username, cts.RoleAdmin, token)
	if errBody != nil {
		return newToken, errBody
	}

	return repo.Delete(id)
}

func UpdateUser(username, target, role string, token models.Token) (any, *models.ErrorBody) {
	newToken, errBody := preprocessUser(username, cts.RoleAdmin, token)
	if errBody != nil {
		return newToken, errBody
	}

	err := kk.UpdateUser(target, role)
	if err != nil {
		return newToken, &models.ErrorBody{Error: err.Error(), ErrorCode: http.StatusNotFound}
	}
	return newToken, nil
}

func CreateUser(username, password string) (any, *models.ErrorBody) {
	kk.RegisterUser(username, password, cts.RoleCardOwner)
	return repo.CreateProfileRequest(username)
}

func preprocessUser(username, requiredRole string, token models.Token) (*models.Token, *models.ErrorBody) {
	newToken, err := kk.CheckToken(token)
	if err != nil {
		return newToken, &models.ErrorBody{Error: err.Error(), ErrorCode: http.StatusBadRequest}
	}

	hasRole, err := kk.CheckRole(username, requiredRole)
	if err != nil {
		return newToken, &models.ErrorBody{Error: err.Error(), ErrorCode: http.StatusBadRequest}
	} else if !hasRole {
		return newToken, &models.ErrorBody{Error: err.Error(), ErrorCode: http.StatusForbidden}
	}

	return newToken, nil
}