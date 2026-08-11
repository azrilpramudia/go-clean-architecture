package http

import (
	"encoding/json"
	"net/http"

	"github.com/azrilpramudia/go-clean-architecture/internal/model"
	"github.com/azrilpramudia/go-clean-architecture/internal/usecase"
)

type UserHandler struct {
	Usecase *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{Usecase: uc}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	request := new(model.RegisterUserRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.Usecase.Register(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}