package user

import (
	"fmt"
	"net/http"

	"github.com/cityzenK/com/service/auth"
	"github.com/cityzenK/com/types"
	"github.com/cityzenK/com/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	//Get JSON
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// validate payload

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	//Check if recurse exists
	_, err := h.store.GetUserByEmail(payload.Email)

	if err == nil {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			fmt.Errorf("user with email %s alredy exists", payload.Email),
		)
		return
	}

	hashPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	err = h.store.CreateUser(types.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
		UserName: hashPassword,
	})

	if err != nil {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
