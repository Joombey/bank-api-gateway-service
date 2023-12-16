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
	newToken, errBody := preprocessUser(username, kk.RoleCardOwner, token)
	if errBody != nil {
		return models.BaseHTTPModel{Err: *errBody}
	}

	serviceResponse, errFromResponse := api.TransferRequest(from, to, value)
	return models.BaseHTTPModel{Body: serviceResponse, Token: *newToken, Err: errFromResponse}
}

func CreateUser(username, password string) models.BaseHTTPModel {
	kk.RegisterUser(username, password, kk.RoleUser)
	response, errBody := api.CreateProfile(username)
	
	token, err := kk.LoginUser(username, password)
	if err != nil {
		panic(err)
	}
	
	kk.UpdateUser(username, kk.RoleCardOwner)
	return models.BaseHTTPModel{Body: response, Err: errBody, Token: token}
}

func BlockUser(id int, username, target string, token models.Token) models.BaseHTTPModel {
	Token, errBody := preprocessUser(username, kk.RoleAdmin, token)
	if errBody != nil {
		if Token == nil {
			return models.BaseHTTPModel{Err: *errBody}
		} else {
			return models.BaseHTTPModel{Token: *Token, Err: *errBody}
		}
	}

	deleteErr := kk.DeleteUser(target)
	if deleteErr != nil {
		panic(deleteErr.Error())
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
	if username == "" {
		return nil, &models.ErrorBody{ Error: "username must not be empty", ErrorCode: http.StatusBadRequest}
	}
	
	newToken, err := kk.CheckToken(token)
	if err != nil {
		return newToken, &models.ErrorBody{Error: err.Error(), ErrorCode: http.StatusBadRequest}
	}

	hasRole, err := kk.CheckRole(username, requiredRole)
	if err != nil {
		return nil, &models.ErrorBody{Error: err.Error(), ErrorCode: http.StatusBadRequest}
	} else if !hasRole {
		return nil, &models.ErrorBody{Error: "User has not acces for that type of action", ErrorCode: http.StatusForbidden}
	}

	println(hasRole)

	return newToken, nil
}
