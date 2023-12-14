package handlers

import (
	"net/http"

	api "farukh.go/api-gateway/api"
	kk "farukh.go/api-gateway/keycloak"
	"farukh.go/api-gateway/models"
)

func GetProfileHandler(id int) models.BaseHTTPModel {
	response, err := api.GetProfile(id)
	return models.BaseHTTPModel{Body: response, Err: err}
}

func Transfer(username string, token models.Token, from int, to int, value float32) models.BaseHTTPModel {
	Token, errBody := preprocessUser(username, kk.RoleCardOwner, token)
	if errBody != nil {
		return models.BaseHTTPModel{Token: *Token, Err: *errBody}
	}

	serviceResponse, errFromResponse := api.TransferRequest(from, to, value)
	return models.BaseHTTPModel{Body: serviceResponse, Token: *Token, Err: errFromResponse}
}

func CreateUser(username, password string) models.BaseHTTPModel {
	kk.RegisterUser(username, password, kk.RoleCardOwner)
	response, err := api.CreateProfile(username)
	return models.BaseHTTPModel{Body: response, Err: err}
}

func BlockUser(id int, username, target string, token models.Token) models.BaseHTTPModel {
	Token, errBody := preprocessUser(username, kk.RoleAdmin, token)
	if errBody != nil {
		return models.BaseHTTPModel{Token: *Token, Err: *errBody}
	}

	responseModel, err := api.BlockUser(id)
	return models.BaseHTTPModel{Body: responseModel, Token: *Token, Err: err}
}

func LoadMoney(request models.InsertRequest) models.BaseHTTPModel{
	response, err := api.LoadMoney(request)
	return models.BaseHTTPModel {
		Body: response,
		Err: err,
	}
}

func Login(username, password string) (models.Token, error) {
	return kk.LoginUser(username, password)
}

func UpdateUser(username, target, role string, token models.Token) models.BaseHTTPModel {
	Token, errBody := preprocessUser(username, kk.RoleAdmin, token)
	if errBody != nil {
		return models.BaseHTTPModel{Token: *Token, Err: *errBody}
	}

	err := kk.UpdateUser(target, role)
	if err != nil {
		return models.BaseHTTPModel{Token: *Token, Err: models.ErrorBody{Error: err.Error(), ErrorCode: http.StatusNotFound}}
	}
	return models.BaseHTTPModel{Token: *Token}
}

func preprocessUser(username, requiredRole string, token models.Token) (*models.Token, *models.ErrorBody) {
	Token, err := kk.CheckToken(token)
	if err != nil {
		return Token, &models.ErrorBody{Error: err.Error(), ErrorCode: http.StatusBadRequest}
	}

	hasRole, err := kk.CheckRole(username, requiredRole)
	if err != nil {
		return Token, &models.ErrorBody{Error: err.Error(), ErrorCode: http.StatusBadRequest}
	} else if !hasRole {
		return Token, &models.ErrorBody{Error: err.Error(), ErrorCode: http.StatusForbidden}
	}

	return Token, nil
}
