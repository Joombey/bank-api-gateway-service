package models

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Jwt      Token  `json:"token"`
}

type DeleteUserRequest struct {
	Username string `json:"username"`
	Caller   string `json:"caller"`
	Jwt      Token  `json:"token"`
}

type ValueResponse struct {
	CardNumber int     `json:"card_number"`
	Value      float32 `json:"value"`
	Jwt        Token   `json:"token"`
}

type ProfileResponse struct {
	Id         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	CardNumber int    `db:"card_number" json:"card_number"`
	Value      string `db:"password" json:"password"`
}

type TransferDTO struct {
	From  int     `json:"from"`
	To    int     `json:"to"`
	Value float32 `json:"value"`
	Jwt   Token   `json:"token"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TransferResponse struct {
	NewValues []ValueResponse `json:"new_values"`
	NewToken  Token           `json:"new_token"`
	Error     string          `json:"error"`
	ErrorCode int
}

type ReponseFrameWithError struct {
	Body     any       `json:"body"`
	NewToken Token     `json:"new_token"`
	Err      ErrorBody `json:"err"`
}

type Profile struct {
	Id         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	CardNumber int    `db:"card_number" json:"card_number"`
}

type ErrorBody struct {
	Error     string `json:"error"`
	ErrorCode int
}
