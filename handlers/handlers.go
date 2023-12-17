package handlers

import (
	"context"
	"net/http"

	api "farukh.go/api-gateway/api"
	kk "farukh.go/api-gateway/keycloak"
	"farukh.go/api-gateway/models"
	tr "farukh.go/api-gateway/tracing"
)

func GetProfileHandler(id int, ctx context.Context) models.BaseHTTPModel {
	_, span := tr.Tracer.Start(ctx, "GET_PROFILE")
	defer span.End()
	
	span.AddEvent("Getting prfile from API")
	response, err := api.GetProfile(id)
	return models.BaseHTTPModel{Body: response, Err: err}
}

func Transfer(username string, token models.Token, from int, to int, value float32, ctx context.Context) models.BaseHTTPModel {
	ctx, span := tr.Tracer.Start(ctx, "Transfer")
	defer span.End()

	span.AddEvent("Checking token and role")
	ctx, prSpan := tr.Tracer.Start(ctx, "Preprocess")
	defer prSpan.End()
	
	newToken, errBody := preprocessUser(username, kk.RoleCardOwner, token)
	if errBody != nil {
		return models.BaseHTTPModel{Err: *errBody}
	}

	span.AddEvent("transfer")
	_, transferReq := tr.Tracer.Start(ctx, "transfer req")
	defer transferReq.End()
	
	serviceResponse, errFromResponse := api.TransferRequest(from, to, value)
	return models.BaseHTTPModel{Body: serviceResponse, Token: *newToken, Err: errFromResponse}
}

func CreateUser(username, password string, ctx context.Context) models.BaseHTTPModel {
	ctx, span := tr.Tracer.Start(ctx, "creating user")
	defer span.End()
	
	span.AddEvent("creating in keycloak")
	ctx, kkRegisterSpan := tr.Tracer.Start(ctx, "registering user")
	defer kkRegisterSpan.End()
	
	_, err := kk.RegisterUser(username, password, kk.RoleUser)
	if err != nil {
		return models.BaseHTTPModel{Err: models.ErrorBody{ Error: err.Error(), ErrorCode: http.StatusConflict} }
	}
	
	span.AddEvent("creating in profile service")
	ctx, createProfileSpan := tr.Tracer.Start(ctx, "registering user")
	defer createProfileSpan.End()

	response, errBody := api.CreateProfile(username)
	
	
	span.AddEvent("login user")
	ctx, loginSpan := tr.Tracer.Start(ctx, "login")
	defer loginSpan.End()

	token, err := kk.LoginUser(username, password)
	if err != nil {
		panic(err)
	}
	
	span.AddEvent("update user role to card owner")
	_, update := tr.Tracer.Start(ctx, "update role")
	defer update.End()

	kk.UpdateUser(username, kk.RoleCardOwner)
	return models.BaseHTTPModel{Body: response, Err: errBody, Token: token}
}

func BlockUser(id int, username, target string, token models.Token, ctx context.Context) models.BaseHTTPModel {
	ctx, span := tr.Tracer.Start(ctx, "blocking")
	defer span.End()
	
	span.AddEvent("preprocess")
	ctx, preSpan := tr.Tracer.Start(ctx, "update role")
	defer preSpan.End()
	
	Token, errBody := preprocessUser(username, kk.RoleAdmin, token)
	if errBody != nil {
		if Token == nil {
			return models.BaseHTTPModel{Err: *errBody}
		} else {
			return models.BaseHTTPModel{Token: *Token, Err: *errBody}
		}
	}

	span.AddEvent("deleting from keycloak")
	ctx, delete := tr.Tracer.Start(ctx, "delete")
	defer delete.End()

	deleteErr := kk.DeleteUser(target)
	if deleteErr != nil {
		panic(deleteErr.Error())
	}

	span.AddEvent("deleting from services")
	_, block := tr.Tracer.Start(ctx, "block")
	defer block.End()

	responseModel, err := api.BlockUser(id)
	return models.BaseHTTPModel{Body: responseModel, Token: *Token, Err: err}
}

func LoadMoney(request models.InsertRequest, ctx context.Context) models.BaseHTTPModel{
	_, span := tr.Tracer.Start(ctx, "load money")
	defer span.End()

	response, err := api.LoadMoney(request)
	return models.BaseHTTPModel {
		Body: response,
		Err: err,
	}
}

func Login(username, password string, ctx context.Context) (models.Token, error) {
	_, span := tr.Tracer.Start(ctx, "login")
	defer span.End()

	span.AddEvent("loggin in keycloak")
	return kk.LoginUser(username, password)
}

func UpdateUser(username, target, role string, token models.Token, ctx context.Context) models.BaseHTTPModel {
	ctx, span := tr.Tracer.Start(ctx, "update user role")
	defer span.End()

	span.AddEvent("preprocess in keycloak, role of caller")
	ctx, preprocess := tr.Tracer.Start(ctx, "preprocess")
	defer preprocess.End()

	Token, errBody := preprocessUser(username, kk.RoleAdmin, token)
	if errBody != nil {
		return models.BaseHTTPModel{Token: *Token, Err: *errBody}
	}

	span.AddEvent("update in keycloak")
	_, updateRole := tr.Tracer.Start(ctx, "updateRole")
	defer updateRole.End()

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

	return newToken, nil
}
