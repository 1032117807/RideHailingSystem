package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"ridehailing/backend/internal/model"
	"ridehailing/backend/internal/repository"
)

type UpdateProfileInput struct {
	Nickname string
	Avatar   string
	Email    string
	Gender   string
	Birthday *time.Time
}

type AccountStatus struct {
	Status           string `json:"status"`
	Role             string `json:"role"`
	RealNameVerified bool   `json:"realNameVerified"`
	CanLogin         bool   `json:"canLogin"`
}

type AdminUserSummary struct {
	TotalUsers     int `json:"totalUsers"`
	PassengerCount int `json:"passengerCount"`
	DriverCount    int `json:"driverCount"`
	AdminCount     int `json:"adminCount"`
	ActiveCount    int `json:"activeCount"`
	FrozenCount    int `json:"frozenCount"`
	DisabledCount  int `json:"disabledCount"`
	VerifiedCount  int `json:"verifiedCount"`
}

type AdminUserListFilter struct {
	Keyword string
	Role    string
	Status  string
}

type AdminUserItem struct {
	ID               uint      `json:"id"`
	Phone            string    `json:"phone"`
	Nickname         string    `json:"nickname"`
	Role             string    `json:"role"`
	RealName         string    `json:"realName"`
	RealNameVerified bool      `json:"realNameVerified"`
	Email            string    `json:"email"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type AdminUpdateUserInput struct {
	Role   string
	Status string
}

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(ctx context.Context, userID uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint, input UpdateProfileInput) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if input.Nickname != "" {
		user.Nickname = strings.TrimSpace(input.Nickname)
	}
	if input.Avatar != "" {
		user.Avatar = strings.TrimSpace(input.Avatar)
	}
	if input.Email != "" {
		user.Email = strings.TrimSpace(input.Email)
	}
	if input.Gender != "" {
		user.Gender = strings.TrimSpace(input.Gender)
	}
	if input.Birthday != nil {
		user.Birthday = input.Birthday
	}
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) VerifyRealName(ctx context.Context, userID uint, realName, idCard string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	realName = strings.TrimSpace(realName)
	idCard = strings.TrimSpace(idCard)
	if realName == "" || idCard == "" {
		return nil, errors.New("realName and idCard are required")
	}

	user.RealName = realName
	user.IDCard = idCard
	user.RealNameVerified = true

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) SwitchRole(ctx context.Context, userID uint, targetRole string) (*model.User, error) {
	return nil, errors.New("role switching is disabled; use a separate account for each role")
}

func (s *UserService) GetAccountStatus(ctx context.Context, userID uint) (*AccountStatus, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &AccountStatus{
		Status:           user.Status,
		Role:             user.Role,
		RealNameVerified: user.RealNameVerified,
		CanLogin:         user.Status == model.UserStatusActive,
	}, nil
}

func (s *UserService) ListUsersForAdmin(ctx context.Context, currentUserID uint, currentUserRole string, filter AdminUserListFilter) ([]*AdminUserItem, error) {
	if currentUserRole != model.RoleAdmin {
		return nil, errors.New("only admin can view user list")
	}

	users, err := s.userRepo.ListForAdmin(ctx, repository.AdminUserListFilter{
		Keyword: strings.TrimSpace(filter.Keyword),
		Role:    normalizeRole(filter.Role),
		Status:  normalizeUserStatus(filter.Status),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*AdminUserItem, 0, len(users))
	for _, user := range users {
		if user == nil {
			continue
		}
		result = append(result, toAdminUserItem(user))
	}
	return result, nil
}

func (s *UserService) GetAdminUserSummary(ctx context.Context, currentUserID uint, currentUserRole string) (*AdminUserSummary, error) {
	if currentUserRole != model.RoleAdmin {
		return nil, errors.New("only admin can view user summary")
	}

	summary, err := s.userRepo.CountSummary(ctx)
	if err != nil {
		return nil, err
	}

	return &AdminUserSummary{
		TotalUsers:     summary.TotalUsers,
		PassengerCount: summary.PassengerCount,
		DriverCount:    summary.DriverCount,
		AdminCount:     summary.AdminCount,
		ActiveCount:    summary.ActiveCount,
		FrozenCount:    summary.FrozenCount,
		DisabledCount:  summary.DisabledCount,
		VerifiedCount:  summary.VerifiedCount,
	}, nil
}

func (s *UserService) UpdateUserByAdmin(ctx context.Context, currentUserID uint, currentUserRole string, targetUserID uint, input AdminUpdateUserInput) (*AdminUserItem, error) {
	if currentUserRole != model.RoleAdmin {
		return nil, errors.New("only admin can update user")
	}
	if targetUserID == 0 {
		return nil, errors.New("invalid user id")
	}
	if currentUserID == targetUserID {
		return nil, errors.New("current admin cannot modify own role or status")
	}

	user, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	nextRole := user.Role
	if strings.TrimSpace(input.Role) != "" {
		nextRole = normalizeRole(input.Role)
		if nextRole != model.RolePassenger && nextRole != model.RoleDriver && nextRole != model.RoleAdmin {
			return nil, errors.New("invalid role")
		}
	}

	nextStatus := user.Status
	if strings.TrimSpace(input.Status) != "" {
		nextStatus = normalizeUserStatus(input.Status)
		if nextStatus != model.UserStatusActive && nextStatus != model.UserStatusFrozen && nextStatus != model.UserStatusDisabled {
			return nil, errors.New("invalid status")
		}
	}

	if nextRole == user.Role && nextStatus == user.Status {
		return toAdminUserItem(user), nil
	}

	if user.Role == model.RoleAdmin && user.Status == model.UserStatusActive &&
		(nextRole != model.RoleAdmin || nextStatus != model.UserStatusActive) {
		activeAdmins, err := s.userRepo.CountActiveAdmins(ctx)
		if err != nil {
			return nil, err
		}
		if activeAdmins <= 1 {
			return nil, errors.New("at least one active admin account must remain")
		}
	}

	user.Role = nextRole
	user.DefaultRole = nextRole
	user.Status = nextStatus

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return toAdminUserItem(user), nil
}

func toAdminUserItem(user *model.User) *AdminUserItem {
	if user == nil {
		return nil
	}

	return &AdminUserItem{
		ID:               user.ID,
		Phone:            user.Phone,
		Nickname:         user.Nickname,
		Role:             user.Role,
		RealName:         user.RealName,
		RealNameVerified: user.RealNameVerified,
		Email:            user.Email,
		Status:           user.Status,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}
}

func normalizeUserStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}
