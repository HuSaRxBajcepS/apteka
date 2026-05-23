package handlers

import (
	"encoding/json"
	"net/http"

	"apteka/internal/auth"
	"apteka/internal/services"
)

type AuthHandler struct {
	Auth      *services.AuthService
	JWTSecret string
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", 400)
		return
	}

	role, email, err := h.Auth.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	var userID int
	err = h.Auth.DB.QueryRow(`SELECT id FROM users WHERE email=$1`, email).Scan(&userID)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	token, err := auth.GenerateToken(userID, role, h.JWTSecret)
	if err != nil {
		http.Error(w, "token error", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"token": token,
		"role":  role,
	},
	)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid body", 400)
		return
	}

	err = h.Auth.CreateUser(req.FullName, req.Email, req.Password, req.Role)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
}
