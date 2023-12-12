package profile

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

var client = http.Client{Timeout: time.Duration(10) * time.Second}

type ProfileService struct {
	client http.Client
}

func New() *ProfileService {
	return &ProfileService{client: client}
}

func (p *ProfileService) CreateProfile() {
	
}

func DecodeJson(reader io.Reader, data *any) {
	json.NewDecoder(reader).Decode(data)
}