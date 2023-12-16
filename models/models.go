package models

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type DeleteUserRequest struct {
	Username string `json:"username"`
	Caller   string `json:"caller"`
}

type ValueResponse struct {
	CardNumber int     `json:"card_number"`
	Value      float32 `json:"value"`
}

type ProfileResponse struct {
	Id         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	CardNumber int    `db:"card_number" json:"card_number"`
	Value      string `db:"password" json:"password"`
}

type TransferDTO struct {
	Username string  `json:"username,omitempty"`
	From     int     `json:"from"`
	To       int     `json:"to"`
	Value    float32 `json:"value"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type BaseHTTPModel struct {
	Body  any       `json:"body,omitempty"`
	Token Token     `json:"new_token,omitempty"`
	Err   ErrorBody `json:"err",omitempty`
}

type Profile struct {
	Id         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	CardNumber int    `db:"card_number" json:"card_number"`
}

type ErrorBody struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"-"`
}

type InsertRequest struct {
	CardNumber int     `json:"card_number"`
	Value      float32 `json:"value"`
}
