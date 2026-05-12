package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"

	"ridehailing/backend/internal/model"
)

type TripRepository interface {
	CreateWithStops(ctx context.Context, trip *model.Trip, stops []model.TripStop) error
	ListByDriverID(ctx context.Context, driverID uint) ([]*model.Trip, error)
	GetByID(ctx context.Context, id uint) (*model.Trip, error)
	SearchPublishedTrips(ctx context.Context, startCity, endCity string, departureStart, departureEnd time.Time) ([]*model.Trip, error)
	ListPublishedTripsByDate(ctx context.Context, departureStart, departureEnd time.Time) ([]*model.Trip, error)
}

type GormTripRepository struct {
	db *gorm.DB
}

func NewGormTripRepository(db *gorm.DB) *GormTripRepository {
	return &GormTripRepository{db: db}
}

// 创建行程和停靠点的事务方法
func (r *GormTripRepository) CreateWithStops(ctx context.Context, trip *model.Trip, stops []model.TripStop) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(trip).Error; err != nil {
			return err
		}
		if len(stops) == 0 {
			return nil
		}

		for i := range stops {
			stops[i].TripID = trip.ID
		}

		if err := tx.Create(&stops).Error; err != nil {
			return err
		}

		trip.Stops = stops
		return nil
	})
}

// 根据司机ID查询行程列表
func (r *GormTripRepository) ListByDriverID(ctx context.Context, driverID uint) ([]*model.Trip, error) {
	var trips []*model.Trip
	err := r.db.WithContext(ctx).
		Preload("Stops", func(db *gorm.DB) *gorm.DB {
			return db.Order("stop_order asc")
		}).
		Where("driver_id = ?", driverID).
		Order("departure_time asc").
		Find(&trips).Error
	if err != nil {
		return nil, err
	}
	return trips, nil
}

// 根据ID查询行程详情
func (r *GormTripRepository) GetByID(ctx context.Context, id uint) (*model.Trip, error) {
	var trip model.Trip
	err := r.db.WithContext(ctx).
		Preload("Stops", func(db *gorm.DB) *gorm.DB {
			return db.Order("stop_order asc")
		}).
		First(&trip, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

// 搜索已发布的行程
func (r *GormTripRepository) SearchPublishedTrips(ctx context.Context, startCity, endCity string, departureStart, departureEnd time.Time) ([]*model.Trip, error) {
	var trips []*model.Trip
	err := r.db.WithContext(ctx).
		Preload("Stops", func(db *gorm.DB) *gorm.DB {
			return db.Order("stop_order asc")
		}).
		Where("start_city = ? AND end_city = ?", startCity, endCity).
		Where("departure_time >= ? AND departure_time < ?", departureStart, departureEnd).
		Where("status = ?", model.TripStatusPublished).
		Where("seat_available > 0").
		Order("departure_time asc").
		Find(&trips).Error
	if err != nil {
		return nil, err
	}
	return trips, nil
}

// 查询指定日期内所有可售的已发布班次，供中转拼接使用。
func (r *GormTripRepository) ListPublishedTripsByDate(ctx context.Context, departureStart, departureEnd time.Time) ([]*model.Trip, error) {
	var trips []*model.Trip
	err := r.db.WithContext(ctx).
		Preload("Stops", func(db *gorm.DB) *gorm.DB {
			return db.Order("stop_order asc")
		}).
		Where("departure_time >= ? AND departure_time < ?", departureStart, departureEnd).
		Where("status = ?", model.TripStatusPublished).
		Where("seat_available > 0").
		Order("departure_time asc").
		Find(&trips).Error
	if err != nil {
		return nil, err
	}
	return trips, nil
}
