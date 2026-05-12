package handler

import (
	"net/http"
	"strconv"
	"strings"

	"ridehailing/backend/internal/pkg/middleware"
	"ridehailing/backend/internal/pkg/response"
	"ridehailing/backend/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type updateProfileRequest struct {
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	Birthday    string `json:"birthday"`
	DefaultRole string `json:"defaultRole"`
}

type verifyRealNameRequest struct {
	RealName string `json:"realName"`
	IDCard   string `json:"idCard"`
}

type switchRoleRequest struct {
	TargetRole string `json:"targetRole"`
}

type adminUpdateUserRequest struct {
	Role   string `json:"role"`
	Status string `json:"status"`
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.userService.GetProfile(r.Context(), currentUser.ID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.Success(w, user)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req updateProfileRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	birthday, err := parseBirthday(req.Birthday)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.UpdateProfile(r.Context(), currentUser.ID, service.UpdateProfileInput{
		Nickname:    req.Nickname,
		Avatar:      req.Avatar,
		Email:       req.Email,
		Gender:      req.Gender,
		Birthday:    birthday,
		DefaultRole: req.DefaultRole,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, user)
}

func (h *UserHandler) VerifyRealName(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req verifyRealNameRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.VerifyRealName(r.Context(), currentUser.ID, req.RealName, req.IDCard)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, user)
}

func (h *UserHandler) SwitchRole(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req switchRoleRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.SwitchRole(r.Context(), currentUser.ID, req.TargetRole)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, user)
}

func (h *UserHandler) GetAccountStatus(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	status, err := h.userService.GetAccountStatus(r.Context(), currentUser.ID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.Success(w, status)
}

func (h *UserHandler) ListUsersForAdmin(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	users, err := h.userService.ListUsersForAdmin(r.Context(), currentUser.ID, currentUser.Role, service.AdminUserListFilter{
		Keyword: strings.TrimSpace(r.URL.Query().Get("keyword")),
		Role:    strings.TrimSpace(r.URL.Query().Get("role")),
		Status:  strings.TrimSpace(r.URL.Query().Get("status")),
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, users)
}

func (h *UserHandler) GetAdminUserSummary(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	summary, err := h.userService.GetAdminUserSummary(r.Context(), currentUser.ID, currentUser.Role)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, summary)
}

func (h *UserHandler) UpdateUserByAdmin(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := middleware.CurrentUser(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := strconv.ParseUint(strings.TrimSpace(r.PathValue("userId")), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid userId")
		return
	}

	var req adminUpdateUserRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.UpdateUserByAdmin(r.Context(), currentUser.ID, currentUser.Role, uint(userID), service.AdminUpdateUserInput{
		Role:   req.Role,
		Status: req.Status,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, user)
}
