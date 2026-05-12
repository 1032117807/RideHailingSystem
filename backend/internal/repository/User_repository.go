package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"ridehailing/backend/internal/model"
)

type AdminUserListFilter struct {
	Keyword string
	Role    string
	Status  string
	Limit   int
}

type UserSummaryCount struct {
	TotalUsers     int
	PassengerCount int
	DriverCount    int
	AdminCount     int
	ActiveCount    int
	FrozenCount    int
	DisabledCount  int
	VerifiedCount  int
}

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	ListForAdmin(ctx context.Context, filter AdminUserListFilter) ([]*model.User, error)
	CountSummary(ctx context.Context) (*UserSummaryCount, error)
	CountActiveAdmins(ctx context.Context) (int, error)
	Update(ctx context.Context, user *model.User) error
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *GormUserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) ListForAdmin(ctx context.Context, filter AdminUserListFilter) ([]*model.User, error) {
	var users []*model.User
	query := r.db.WithContext(ctx).Order("created_at desc, id desc")

	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			r.db.Where("phone LIKE ?", like).
				Or("nickname LIKE ?", like).
				Or("email LIKE ?", like).
				Or("real_name LIKE ?", like),
		)
	}

	if role := strings.TrimSpace(filter.Role); role != "" {
		query = query.Where("role = ?", role)
	}

	if status := strings.TrimSpace(filter.Status); status != "" {
		query = query.Where("status = ?", status)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *GormUserRepository) CountSummary(ctx context.Context) (*UserSummaryCount, error) {
	countUsers := func(query *gorm.DB) (int, error) {
		var count int64
		if err := query.Count(&count).Error; err != nil {
			return 0, err
		}
		return int(count), nil
	}

	base := r.db.WithContext(ctx).Model(&model.User{})

	total, err := countUsers(base)
	if err != nil {
		return nil, err
	}
	passenger, err := countUsers(base.Where("role = ?", model.RolePassenger))
	if err != nil {
		return nil, err
	}
	driver, err := countUsers(base.Where("role = ?", model.RoleDriver))
	if err != nil {
		return nil, err
	}
	admin, err := countUsers(base.Where("role = ?", model.RoleAdmin))
	if err != nil {
		return nil, err
	}
	active, err := countUsers(base.Where("status = ?", model.UserStatusActive))
	if err != nil {
		return nil, err
	}
	frozen, err := countUsers(base.Where("status = ?", model.UserStatusFrozen))
	if err != nil {
		return nil, err
	}
	disabled, err := countUsers(base.Where("status = ?", model.UserStatusDisabled))
	if err != nil {
		return nil, err
	}
	verified, err := countUsers(base.Where("real_name_verified = ?", true))
	if err != nil {
		return nil, err
	}

	return &UserSummaryCount{
		TotalUsers:     total,
		PassengerCount: passenger,
		DriverCount:    driver,
		AdminCount:     admin,
		ActiveCount:    active,
		FrozenCount:    frozen,
		DisabledCount:  disabled,
		VerifiedCount:  verified,
	}, nil
}

func (r *GormUserRepository) CountActiveAdmins(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("role = ? AND status = ?", model.RoleAdmin, model.UserStatusActive).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *GormUserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
