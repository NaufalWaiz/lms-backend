package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/service"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

type UserHandler struct {
	UserService *service.UserService
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Role == "" {
		// Set default role jika tidak diberikan
		request.Role = "Murid"
	}

	err := h.AuthService.Register(request.Username, request.Password, request.Role)
	if err != nil {
		fmt.Printf("Failed to register user: %v", err) // Log the error for debugging
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Register Account successfully",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.Login(request.Username, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) GetRoleCounts(w http.ResponseWriter, r *http.Request) {
    // Panggil service untuk menghitung pengguna berdasarkan role
	roleCounts, err := h.UserService.CountUsersByRole()
    if err != nil {
        http.Error(w, "Failed to retrieve role counts: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Kirim response sebagai JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(roleCounts)
}