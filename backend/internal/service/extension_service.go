package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"ridehailing/backend/internal/model"
	"ridehailing/backend/internal/repository"
)

type PassengerService struct {
	passengerRepo repository.PassengerRepository
}

func NewPassengerService(passengerRepo repository.PassengerRepository) *PassengerService {
	return &PassengerService{passengerRepo: passengerRepo}
}

func (s *PassengerService) Create(ctx context.Context, userID uint, role string, passenger *model.Passenger) (*model.Passenger, error) {
	if role != model.RolePassenger {
		return nil, errors.New("only passenger can manage passengers")
	}
	passenger.UserID = userID
	passenger.Name = strings.TrimSpace(passenger.Name)
	passenger.IDCard = strings.TrimSpace(passenger.IDCard)
	passenger.Phone = strings.TrimSpace(passenger.Phone)
	if passenger.Name == "" || passenger.IDCard == "" || passenger.Phone == "" {
		return nil, errors.New("name, idCard and phone are required")
	}
	if err := s.passengerRepo.Create(ctx, passenger); err != nil {
		return nil, err
	}
	return passenger, nil
}

func (s *PassengerService) List(ctx context.Context, userID uint, role string) ([]*model.Passenger, error) {
	if role != model.RolePassenger {
		return nil, errors.New("only passenger can view passengers")
	}
	return s.passengerRepo.ListByUserID(ctx, userID)
}

func (s *PassengerService) Delete(ctx context.Context, userID uint, role string, passengerID uint) error {
	if role != model.RolePassenger {
		return errors.New("only passenger can delete passengers")
	}
	return s.passengerRepo.DeleteByUser(ctx, userID, passengerID)
}

type PriceAlertService struct {
	priceAlertRepo repository.PriceAlertRepository
}

func NewPriceAlertService(priceAlertRepo repository.PriceAlertRepository) *PriceAlertService {
	return &PriceAlertService{priceAlertRepo: priceAlertRepo}
}

func (s *PriceAlertService) Create(ctx context.Context, userID uint, role string, alert *model.PriceAlert) (*model.PriceAlert, error) {
	if role != model.RolePassenger {
		return nil, errors.New("only passenger can create price alert")
	}
	alert.UserID = userID
	alert.StartCity = strings.TrimSpace(alert.StartCity)
	alert.EndCity = strings.TrimSpace(alert.EndCity)
	alert.Status = model.PriceAlertStatusActive
	if alert.StartCity == "" || alert.EndCity == "" || alert.TargetPriceCent <= 0 {
		return nil, errors.New("startCity, endCity and targetPriceCent are required")
	}
	if alert.StartDate.IsZero() {
		alert.StartDate = time.Now()
	}
	if alert.EndDate.IsZero() || alert.EndDate.Before(alert.StartDate) {
		alert.EndDate = alert.StartDate.AddDate(0, 0, 7)
	}
	if err := s.priceAlertRepo.Create(ctx, alert); err != nil {
		return nil, err
	}
	return alert, nil
}

func (s *PriceAlertService) List(ctx context.Context, userID uint, role string) ([]*model.PriceAlert, error) {
	if role != model.RolePassenger {
		return nil, errors.New("only passenger can view price alerts")
	}
	return s.priceAlertRepo.ListByUserID(ctx, userID)
}

func (s *PriceAlertService) Disable(ctx context.Context, userID uint, role string, alertID uint) error {
	if role != model.RolePassenger {
		return errors.New("only passenger can disable price alert")
	}
	return s.priceAlertRepo.DisableByUser(ctx, userID, alertID)
}

type DriverProfileService struct {
	driverRepo repository.DriverProfileRepository
}

func NewDriverProfileService(driverRepo repository.DriverProfileRepository) *DriverProfileService {
	return &DriverProfileService{driverRepo: driverRepo}
}

func (s *DriverProfileService) UpsertProfile(ctx context.Context, userID uint, role string, profile *model.DriverProfile) (*model.DriverProfile, error) {
	if role != model.RoleDriver {
		return nil, errors.New("only driver can submit driver profile")
	}
	profile.UserID = userID
	if strings.TrimSpace(profile.RealName) == "" || strings.TrimSpace(profile.IDCard) == "" || strings.TrimSpace(profile.LicenseNo) == "" {
		return nil, errors.New("realName, idCard and licenseNo are required")
	}
	profile.Status = model.DriverProfileStatusPending
	if err := s.driverRepo.UpsertProfile(ctx, profile); err != nil {
		return nil, err
	}
	return s.driverRepo.GetByUserID(ctx, userID)
}

func (s *DriverProfileService) GetProfile(ctx context.Context, userID uint, role string) (*model.DriverProfile, error) {
	if role != model.RoleDriver {
		return nil, errors.New("only driver can view driver profile")
	}
	return s.driverRepo.GetByUserID(ctx, userID)
}

func (s *DriverProfileService) CreateVehicle(ctx context.Context, userID uint, role string, vehicle *model.Vehicle) (*model.Vehicle, error) {
	if role != model.RoleDriver {
		return nil, errors.New("only driver can manage vehicles")
	}
	vehicle.DriverID = userID
	if strings.TrimSpace(vehicle.PlateNo) == "" || vehicle.SeatCount <= 0 {
		return nil, errors.New("plateNo and seatCount are required")
	}
	if vehicle.Status == "" {
		vehicle.Status = model.UserStatusActive
	}
	if err := s.driverRepo.CreateVehicle(ctx, vehicle); err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (s *DriverProfileService) ListVehicles(ctx context.Context, userID uint, role string) ([]*model.Vehicle, error) {
	if role != model.RoleDriver {
		return nil, errors.New("only driver can view vehicles")
	}
	return s.driverRepo.ListVehicles(ctx, userID)
}
