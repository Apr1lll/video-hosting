package handlers

import (
	"auth/internal/api/dto"
	"auth/internal/domain"
	"auth/jwt"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	userService *domain.Service
	jwtGen      *jwt.Generator
}

func NewAuthHandler(userService *domain.Service, jwtGen *jwt.Generator) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtGen:      jwtGen,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	accessToken, err := h.jwtGen.GenerateAccessToken(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	refreshToken, _ := h.jwtGen.GenerateRefreshToken(user)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   604800,
	})

	respondWithJSON(w, http.StatusOK, dto.LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   15,
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, dto.ErrorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
