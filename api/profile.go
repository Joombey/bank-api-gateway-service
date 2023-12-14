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
	profileUrl = baseUrl + "profile"
)

func TransferRequest(from int, to int, value float32) ([]models.ValueResponse, models.ErrorBody) {
	requestBody := models.TransferDTO{From: from, To: to, Value: value}
	buffer := prepareJson(requestBody)

	response, err := client.Post(sendUrl, "application/json", buffer)
	if err != nil {
		return nil, models.ErrorBody{Error: err.Error(), ErrorCode: 500}
	}

	modelResponseModel := make([]models.ValueResponse, 2)
	decodeJson(response.Body, &modelResponseModel)

	return modelResponseModel, models.ErrorBody{Error: response.Status, ErrorCode: response.StatusCode}
}

func CreateProfile(username string) (*models.Profile, models.ErrorBody) {
	finalUrl := fmt.Sprintf("%s/%s", createUrl, username)
	response, err := client.Get(finalUrl)
	if err != nil {
		return nil, models.ErrorBody{Error: err.Error(), ErrorCode: 500}
	}

	responseModel := models.Profile{}
	decodeJson(response.Body, &responseModel)
	return &responseModel, models.ErrorBody{Error: response.Status, ErrorCode: response.StatusCode}
}

func GetProfile(id int) (*models.Profile, models.ErrorBody){
	finalUrl := fmt.Sprintf("%s/%d", createUrl, id)
	response, err := client.Get(finalUrl)
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
