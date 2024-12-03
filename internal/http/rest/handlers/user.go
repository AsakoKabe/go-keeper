package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"go-keeper/internal/db/repository/errs"
	"go-keeper/internal/keeper/user"
	"go-keeper/internal/models"
	"go-keeper/internal/utils/jwt"
	"go-keeper/pkg/middleware"
)

// UserHandler Структура для user endpoints
type UserHandler struct {
	userService *user.Service
}

// NewUserHandler Конструктор для UserHandler
func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user *models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, err := h.userService.CreateUser(r.Context(), user)
	if errors.Is(err, errs.ErrLoginAlreadyExist) {
		slog.Error("login already exist")
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		slog.Error("error to add user to db", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenString, err := jwt.BuildJWTString(*userID)
	if err != nil {
		slog.Error("error to create jwt token", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.SetCookie(
		w, &http.Cookie{
			Name:  middleware.CookieName,
			Value: tokenString,
		},
	)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("ok"))
}

func (h *UserHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user *models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	existedUser, err := h.userService.GetUser(r.Context(), user)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		slog.Error("error to check user exist", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if existedUser == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !h.userService.CorrectCredentials(user, existedUser) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString, err := jwt.BuildJWTString(existedUser.ID)
	if err != nil {
		slog.Error("error to create jwt token", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.SetCookie(
		w, &http.Cookie{
			Name:  middleware.CookieName,
			Value: tokenString,
		},
	)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
