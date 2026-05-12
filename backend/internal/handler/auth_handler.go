package handler

import (
	"net/http"

	"ridehailing/backend/internal/pkg/middleware"
	"ridehailing/backend/internal/pkg/response"
	"ridehailing/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type sendEmailCodeRequest struct {
	Email string `json:"email"`
	Scene string `json:"scene"`
}

type registerRequest struct {
	Role      string `json:"role"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	EmailCode string `json:"emailCode"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
}

type loginByPasswordRequest struct {
	Role     string `json:"role"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type loginByCodeRequest struct {
	Role      string `json:"role"`
	Email     string `json:"email"`
	EmailCode string `json:"emailCode"`
}

func (h *AuthHandler) SendEmailCode(w http.ResponseWriter, r *http.Request) {
	var req sendEmailCodeRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	code, err := h.authService.SendEmailCode(r.Context(), req.Email, req.Scene)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, map[string]interface{}{
		"email": req.Email,
		"scene": req.Scene,
		"code":  code,
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.authService.Register(r.Context(), req.Role, req.Phone, req.Email, req.EmailCode, req.Password, req.Nickname)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, result)
}

func (h *AuthHandler) LoginByPassword(w http.ResponseWriter, r *http.Request) {
	var req loginByPasswordRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.authService.LoginByPassword(r.Context(), req.Role, req.Phone, req.Password)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, result)
}

func (h *AuthHandler) LoginByCode(w http.ResponseWriter, r *http.Request) {
	var req loginByCodeRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.authService.LoginByCode(r.Context(), req.Role, req.Email, req.EmailCode)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, result)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.authService.Logout(r.Context(), currentUser.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, map[string]string{
		"message": "logout success",
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.authService.Me(r.Context(), currentUser.ID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.Success(w, user)
}
