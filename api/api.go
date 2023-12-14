package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"farukh.go/api-gateway/models"
)

var client = http.Client{Timeout: time.Duration(10) * time.Second}

const (
	baseUrl    = "http://profile-app:8080/"
	sendUrl    = baseUrl + "send"
	createUrl  = baseUrl + "create"
	profileUrl = baseUrl + "credentials"
	blockUrl   = baseUrl + "block"
	loadUrl   = "http://bank-app:8080/load-money"
)

func TransferRequest(from int, to int, value float32) ([]models.ValueResponse, models.ErrorBody) {
	requestBody := models.TransferDTO{From: from, To: to, Value: value}
	modelResponseModel := make([]models.ValueResponse, 2)
	err := makeBasePostRequest(sendUrl, requestBody, &modelResponseModel)

	return modelResponseModel, err
}

func LoadMoney(request any) (models.ValueResponse, models.ErrorBody) {
	var resposeModel models.ValueResponse
	err := makeBasePostRequest(loadUrl, request, &resposeModel)
	return resposeModel, err
}

func CreateProfile(username string) (*models.Profile, models.ErrorBody) {
	finalUrl := fmt.Sprintf("%s/%s", createUrl, username)
	return makeBaseProfileGetRequest(finalUrl)
}

func GetProfile(id int) (*models.Profile, models.ErrorBody) {
	finalUrl := fmt.Sprintf("%s/%d", createUrl, id)
	return makeBaseProfileGetRequest(finalUrl)
}

func BlockUser(userID int) (*models.Profile, models.ErrorBody) {
	finalUrl := fmt.Sprintf("%s/%d", createUrl, userID)
	return makeBaseProfileGetRequest(finalUrl)
}

func UploadMoney(cardNumber int, value float32) {

}

func makeBasePostRequest(url string, model any, responseModel any) (models.ErrorBody) {
	response, err := client.Post(url, "application/json", prepareJson(model))
	if err != nil {
		return models.ErrorBody{Error: err.Error(), ErrorCode: 500}
	}
	
	decodeJson(response.Body, responseModel)
	return models.ErrorBody{Error: response.Status, ErrorCode: response.StatusCode}
}


func makeBaseProfileGetRequest(url string) (*models.Profile, models.ErrorBody) {
	response, err := client.Get(url)
	if err != nil {
		return nil, models.ErrorBody{Error: err.Error(), ErrorCode: 500}
	}

	responseModel := models.Profile{}
	decodeJson(response.Body, &responseModel)
	return &responseModel, models.ErrorBody{Error: response.Status, ErrorCode: response.StatusCode}
}

func decodeJson(reader io.Reader, data any) {
	json.NewDecoder(reader).Decode(data)
}

func prepareJson(jsonModel any) *bytes.Buffer {
	data, err := json.Marshal(jsonModel)
	if err != nil {
		panic(err.Error())
	}

	return bytes.NewBuffer(data)
}
