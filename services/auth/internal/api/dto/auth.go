package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
