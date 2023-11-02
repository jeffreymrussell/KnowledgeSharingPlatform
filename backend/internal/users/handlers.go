package users

import (
	"encoding/json"
	"net/http"
	"strings"
)

type UserHandler struct {
	Usecase *UserUsecase
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var dto RegisterUserDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.Usecase.RegisterUser(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var dto LoginUserDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.Usecase.LoginUser(dto.Username, dto.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
func extractToken(authHeader string) string {
	// Split the header on space and check if it follows the pattern "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r.Header.Get("Authorization"))

	confirmation, err := h.Usecase.LogoutUser(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"confirmation": confirmation})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
