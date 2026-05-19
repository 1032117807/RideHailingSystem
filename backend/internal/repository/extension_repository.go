package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"ridehailing/backend/internal/model"
)

type PassengerRepository interface {
	Create(ctx context.Context, passenger *model.Passenger) error
	ListByUserID(ctx context.Context, userID uint) ([]*model.Passenger, error)
	DeleteByUser(ctx context.Context, userID uint, passengerID uint) error
}

type PriceAlertRepository interface {
	Create(ctx context.Context, alert *model.PriceAlert) error
	ListByUserID(ctx context.Context, userID uint) ([]*model.PriceAlert, error)
	DisableByUser(ctx context.Context, userID uint, alertID uint) error
}

type DriverProfileRepository interface {
	UpsertProfile(ctx context.Context, profile *model.DriverProfile) error
	GetByUserID(ctx context.Context, userID uint) (*model.DriverProfile, error)
	CreateVehicle(ctx context.Context, vehicle *model.Vehicle) error
	ListVehicles(ctx context.Context, driverID uint) ([]*model.Vehicle, error)
}

type GormPassengerRepository struct {
	db *gorm.DB
}

func NewGormPassengerRepository(db *gorm.DB) *GormPassengerRepository {
	return &GormPassengerRepository{db: db}
}

func (r *GormPassengerRepository) Create(ctx context.Context, passenger *model.Passenger) error {
	return r.db.WithContext(ctx).Create(passenger).Error
}

func (r *GormPassengerRepository) ListByUserID(ctx context.Context, userID uint) ([]*model.Passenger, error) {
	var items []*model.Passenger
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&items).Error
	return items, err
}

func (r *GormPassengerRepository) DeleteByUser(ctx context.Context, userID uint, passengerID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, passengerID).Delete(&model.Passenger{}).Error
}

type GormPriceAlertRepository struct {
	db *gorm.DB
}

func NewGormPriceAlertRepository(db *gorm.DB) *GormPriceAlertRepository {
	return &GormPriceAlertRepository{db: db}
}

func (r *GormPriceAlertRepository) Create(ctx context.Context, alert *model.PriceAlert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

func (r *GormPriceAlertRepository) ListByUserID(ctx context.Context, userID uint) ([]*model.PriceAlert, error) {
	var items []*model.PriceAlert
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error
	return items, err
}

func (r *GormPriceAlertRepository) DisableByUser(ctx context.Context, userID uint, alertID uint) error {
	return r.db.WithContext(ctx).
		Model(&model.PriceAlert{}).
		Where("user_id = ? AND id = ?", userID, alertID).
		Update("status", model.PriceAlertStatusDisabled).Error
}

type GormDriverProfileRepository struct {
	db *gorm.DB
}

func NewGormDriverProfileRepository(db *gorm.DB) *GormDriverProfileRepository {
	return &GormDriverProfileRepository{db: db}
}

func (r *GormDriverProfileRepository) UpsertProfile(ctx context.Context, profile *model.DriverProfile) error {
	var existing model.DriverProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", profile.UserID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.WithContext(ctx).Create(profile).Error
	}
	if err != nil {
		return err
	}
	profile.ID = existing.ID
	return r.db.WithContext(ctx).Save(profile).Error
}

func (r *GormDriverProfileRepository) GetByUserID(ctx context.Context, userID uint) (*model.DriverProfile, error) {
	var profile model.DriverProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &profile, err
}

func (r *GormDriverProfileRepository) CreateVehicle(ctx context.Context, vehicle *model.Vehicle) error {
	return r.db.WithContext(ctx).Create(vehicle).Error
}

func (r *GormDriverProfileRepository) ListVehicles(ctx context.Context, driverID uint) ([]*model.Vehicle, error) {
	var items []*model.Vehicle
	err := r.db.WithContext(ctx).Where("driver_id = ?", driverID).Order("created_at DESC").Find(&items).Error
	return items, err
}
